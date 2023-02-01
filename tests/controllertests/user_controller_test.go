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
	"github.com/mitchellh/mapstructure"

	"github.com/PixDale/sh-code-challenge/api/middlewares"
	"github.com/PixDale/sh-code-challenge/api/models"
	"github.com/PixDale/sh-code-challenge/api/responses"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	Authenticate()
	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"Pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			name:         "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "email already taken",
		},
		{
			inputJSON:    `{"name":"Pet", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "name already taken",
		},
		{
			inputJSON:    `{"name":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    `{"name": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required name",
		},
		{
			inputJSON:    `{"name": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required email",
		},
		{
			inputJSON:    `{"name": "Kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "required password",
		},
	}

	for _, v := range samples {
		app := fiber.New()
		app.Post("/users", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.CreateUser)

		req, err := http.NewRequestWithContext(context.Background(), "POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
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
		if v.statusCode == 201 {
			responseUser := (*responseStruct.Data)["data"].(map[string]interface{})
			usr := models.User{}
			usr.Name = responseUser["name"].(string)
			usr.Email = responseUser["email"].(string)
			assert.Equal(t, usr.Name, v.name)
			assert.Equal(t, usr.Email, v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			responseMessage := (*responseStruct.Data)["data"].(string)
			assert.Equal(t, responseMessage, v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Get("/users", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.GetUsers)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/users", nil)
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

	var users []models.User
	responseUsers := (*responseStruct.Data)["data"].([]interface{})
	fmt.Println("aaaaaaaaaaa", responseUsers)
	for _, u := range responseUsers {
		usr := models.User{}
		err := mapstructure.Decode(u, &usr)
		if err != nil {
			t.Errorf("failed to read users: %v\n", err.Error())
		}
		users = append(users, usr)
	}
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Status Code")
	utils.AssertEqual(t, 2, len(users), "Amount of users")
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		name         string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			name:       user.Name,
			email:      user.Email,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range userSample {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		bearer := "Bearer " + managerTokenJWT
		req.Header.Add("Authorization", bearer)

		q := req.URL.Query()
		q.Add("id", v.id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := adaptor.FiberHandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Name, responseMap["name"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers() // we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" // Note the password in the database is already hashed, we want unhashed
	}
	// Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		updateName   string
		updateEmail  string
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Grand", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   200,
			updateName:   "Grand",
			updateEmail:  "grand@gmail.com",
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When password field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Woman", "email": "woman@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Man", "email": "man@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Woman", "email": "woman@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to user 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Frank", "email": "kenny@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email already taken",
		},
		{
			// Remember "Kenny Morris" belongs to user 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Kenny Morris", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "name already taken",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "invalid email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required name",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required email",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"name": "Mike", "email": "mike@gmail.com", "password": "password"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequestWithContext(context.Background(), "POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req.Header.Set("Content-Type", "application/json")
		bearer := "Bearer " + managerTokenJWT
		req.Header.Add("Authorization", bearer)

		q := req.URL.Query()
		q.Add("id", v.id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := adaptor.FiberHandlerFunc(server.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.updateName)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() // we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first and log him in
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}
	// Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// User 2 trying to use User 1 token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range userSample {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		bearer := "Bearer " + managerTokenJWT
		req.Header.Add("Authorization", bearer)

		q := req.URL.Query()
		q.Add("id", v.id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := adaptor.FiberHandlerFunc(server.DeleteUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
