package controllertests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gopkg.in/go-playground/assert.v1"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/PixDale/sh-code-challenge/api/middlewares"
	"github.com/PixDale/sh-code-challenge/api/models"
	"github.com/PixDale/sh-code-challenge/api/responses"
)

func TestCreateTask(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}

	Authenticate()

	samples := []struct {
		inputJSON    string
		statusCode   int
		summary      string
		userID       uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    fmt.Sprintf("{\"summary\": \"The title\", \"user_id\": %d}", managerUser.ID),
			statusCode:   201,
			summary:      "The title",
			userID:       managerUser.ID,
			tokenGiven:   managerTokenJWT,
			errorMessage: "",
		},
		{
			// When no token is passed
			inputJSON:    fmt.Sprintf("{\"summary\": \"When no token is passed\", \"user_id\": %d}", managerUser.ID),
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "error",
		},
		{
			// When incorrect token is passed
			inputJSON:  fmt.Sprintf("{\"summary\": \"When incorrect token is passed\", \"user_id\": %d}", managerUser.ID),
			statusCode: 401,

			tokenGiven:   "incorrect token",
			errorMessage: "error",
		},
		{
			inputJSON:    fmt.Sprintf("{\"summary\": \"\", \"user_id\": %d}", managerUser.ID),
			statusCode:   422,
			tokenGiven:   managerTokenJWT,
			errorMessage: "error",
		},
		{
			inputJSON:    `{"summary": "This is an awesome title"}`,
			statusCode:   400,
			tokenGiven:   managerTokenJWT,
			errorMessage: "error",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    fmt.Sprintf("{\"summary\": \"This is an awesome title\", \"user_id\": %d}", managerUser.ID),
			statusCode:   401,
			tokenGiven:   technicianTokenJWT,
			errorMessage: "error",
		},
	}
	for _, v := range samples {
		req, err := http.NewRequestWithContext(context.Background(), "POST", "/tasks", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Content-Type", "application/json")
		bearer := "Bearer " + v.tokenGiven
		req.Header.Add("Authorization", bearer)
		rr := httptest.NewRecorder()
		handler := adaptor.FiberHandlerFunc(server.CreateTask)
		handler.ServeHTTP(rr, req)

		// responseMap := make(map[string]interface{})
		responseStruct := responses.UserResponse{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseStruct)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			responseTask := (*responseStruct.Data)["data"].(map[string]interface{})
			tsk := models.Task{}
			tsk.Summary = responseTask["summary"].(string)
			tsk.UserID = uint32(responseTask["user_id"].(float64))
			assert.Equal(t, tsk.Summary, v.summary)
			assert.Equal(t, tsk.UserID, v.userID) // just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseStruct.Message, v.errorMessage)
		}
	}
}

func TestGetTasks(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}

	Authenticate()
	_, _, err = seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/tasks", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	bearer := "Bearer " + managerTokenJWT
	req.Header.Add("Authorization", bearer)
	rr := httptest.NewRecorder()
	handler := adaptor.FiberHandlerFunc(server.GetTasks)
	handler.ServeHTTP(rr, req)

	responseStruct := responses.UserResponse{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseStruct)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)

	tasks := (*responseStruct.Data)["data"].([]interface{})
	assert.Equal(t, len(tasks), 2)
}

func TestGetTaskByID(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatal(err)
	}
	taskSample := []struct {
		id           string
		statusCode   int
		summary      string
		userID       uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(task.ID)),
			statusCode: 200,
			summary:    task.Summary,
			userID:     task.UserID,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range taskSample {
		app := fiber.New()
		app.Get("/tasks/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.GetTask)

		url := "/tasks/" + v.id
		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		bearer := "Bearer " + managerTokenJWT
		req.Header.Add("Authorization", bearer)

		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("failed to make the request: %v\n", err.Error())
		}

		responseStruct := responses.UserResponse{}
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v\n", err.Error())
		}
		err = resp.Body.Close()
		if err != nil {
			t.Errorf("failed to close body: %v\n", err.Error())
		}
		err = json.Unmarshal(respBodyBytes, &responseStruct)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		utils.AssertEqual(t, v.statusCode, resp.StatusCode, "Status Code")

		if v.statusCode == 200 {
			task.DecryptSummary()
			responseTask := (*responseStruct.Data)["data"].(map[string]interface{})
			tsk := models.Task{}
			tsk.Summary = responseTask["summary"].(string)
			tsk.UserID = uint32(responseTask["user_id"].(float64))

			assert.Equal(t, task.Summary, tsk.Summary)
			assert.Equal(t, task.UserID, tsk.UserID) // the response author id is float64
		}
	}
}

func TestUpdateTask(t *testing.T) {
	var TaskUserEmail, TaskUserPassword string
	var AuthTaskUserID uint32
	var AuthTaskID uint64

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	users, tasks, err := seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		TaskUserEmail = user.Email
		TaskUserPassword = "password" // Note the password in the database is already hashed, we want unhashed
	}
	// Login the user and get the authentication token
	token, err := server.SignIn(TaskUserEmail, TaskUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first task
	for _, task := range tasks {
		if task.ID == 2 {
			continue
		}
		AuthTaskID = task.ID
		AuthTaskUserID = task.UserID
	}

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		summary      string
		userID       uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"summary":"The updated task", "user_id": 1}`,
			statusCode:   200,
			summary:      "The updated task",
			userID:       AuthTaskUserID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"summary":"This is still another title", "user_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "error",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"summary":"This is still another title", "user_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "error",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"summary":"", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "error",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"summary":"This is another title"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "error",
		},
		{
			id:           "unknown",
			statusCode:   401,
			errorMessage: "error",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"summary":"This is still another title", "user_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "error",
		},
	}

	for _, v := range samples {
		app := fiber.New()
		app.Put("/tasks/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.UpdateTask)

		url := "/tasks/" + v.id
		req, err := http.NewRequestWithContext(context.Background(), "PUT", url, bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", v.tokenGiven)

		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("failed to make the request: %v\n", err.Error())
		}

		responseStruct := responses.UserResponse{}
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v\n", err.Error())
		}
		err = resp.Body.Close()
		if err != nil {
			t.Errorf("failed to close body: %v\n", err.Error())
		}
		err = json.Unmarshal(respBodyBytes, &responseStruct)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		utils.AssertEqual(t, v.statusCode, resp.StatusCode, "Status Code")

		if v.statusCode == 200 {
			responseTask := (*responseStruct.Data)["data"].(map[string]interface{})
			tsk := models.Task{}
			tsk.Summary = responseTask["summary"].(string)
			tsk.UserID = uint32(responseTask["user_id"].(float64))

			assert.Equal(t, tsk.Summary, v.summary)
			assert.Equal(t, tsk.UserID, v.userID) // just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseStruct.Message, v.errorMessage)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	var TaskUserEmail, TaskUserPassword string
	var TaskUserID uint32
	var AuthTaskID uint64

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	users, tasks, err := seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}
	// Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		TaskUserEmail = user.Email
		TaskUserPassword = "password" // Note the password in the database is already hashed, we want unhashed
	}
	// Login the user and get the authentication token
	token, err := server.SignIn(TaskUserEmail, TaskUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second task
	for _, task := range tasks {
		if task.ID == 1 {
			continue
		}
		AuthTaskID = task.ID
		TaskUserID = task.UserID
	}
	taskSample := []struct {
		id           string
		userID       uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTaskID)),
			userID:       TaskUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			userID:       TaskUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "error",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			userID:       TaskUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "error",
		},
		{
			id:           "unknown",
			tokenGiven:   tokenString,
			statusCode:   400,
			errorMessage: "error",
		},
		{
			id:           strconv.Itoa(int(1)),
			userID:       1,
			statusCode:   401,
			errorMessage: "error",
		},
	}
	for _, v := range taskSample {
		app := fiber.New()
		app.Delete("/tasks/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.DeleteTask)

		url := "/tasks/" + v.id
		req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Authorization", v.tokenGiven)

		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("failed to make the request: %v\n", err.Error())
		}

		utils.AssertEqual(t, v.statusCode, resp.StatusCode, "Status Code")

		if v.statusCode == 401 && v.errorMessage != "" {
			responseStruct := responses.UserResponse{}
			respBodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v\n", err.Error())
			}
			err = resp.Body.Close()
			if err != nil {
				t.Errorf("failed to close body: %v\n", err.Error())
			}
			err = json.Unmarshal(respBodyBytes, &responseStruct)
			if err != nil {
				log.Fatalf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseStruct.Message, v.errorMessage)
		}
	}
}
