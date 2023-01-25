package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gopkg.in/go-playground/assert.v1"

	"github.com/gofiber/adaptor/v2"

	"github.com/PixDale/sh-code-challenge/api/models"
)

func TestCreateTask(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") // Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "The title",
			content:      "the content",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Summary",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adaptor.FiberHandlerFunc(server.CreateTask))

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) // just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetTasks(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adaptor.FiberHandlerFunc(server.GetTasks))
	handler.ServeHTTP(rr, req)

	var tasks []models.Task
	err = json.Unmarshal([]byte(rr.Body.String()), &tasks)

	assert.Equal(t, rr.Code, http.StatusOK)
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
		title        string
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(task.ID)),
			statusCode: 200,
			content:    task.Summary,
			author_id:  task.UserID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range taskSample {

		req, err := http.NewRequest("GET", "/tasks", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		q := req.URL.Query()
		q.Add("id", v.id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adaptor.FiberHandlerFunc(server.GetTask))
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, task.Summary, responseMap["content"])
			assert.Equal(t, float64(task.UserID), responseMap["author_id"]) // the response author id is float64
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
	// fmt.Printf("this is the auth task: %v\n", AuthTaskID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"The updated task", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			title:        "The updated task",
			content:      "This is the updated content",
			author_id:    AuthTaskUserID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// Note: "Title 2" belongs to task 2, and title must be unique
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"Title 2", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"Awesome title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Summary",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		q := req.URL.Query()
		q.Add("id", v.id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adaptor.FiberHandlerFunc(server.UpdateTask))

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) // just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
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
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthTaskID)),
			author_id:    TaskUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			author_id:    TaskUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			author_id:    TaskUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range taskSample {

		req, _ := http.NewRequest("GET", "/tasks", nil)
		q := req.URL.Query()
		q.Add("id", v.id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adaptor.FiberHandlerFunc(server.DeleteTask))

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
