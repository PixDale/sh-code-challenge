package controllertests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/adaptor/v2"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "password", // Note the password has to be this, not the hashed one from the database
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {
		t.Logf("[TestSignIn] Trying to login with email: %s and password: %s\n", v.email, v.password)
		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		t.Errorf("this is the error: %v", err)
	}

	_, err = seedOneUser()
	if err != nil {
		t.Errorf("this is the error: %v", err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "pet@gmail.com", "password": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "pet@gmail.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "incorrect password",
		},
		{
			inputJSON:    `{"email": "frank@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "incorrect details",
		},
		{
			inputJSON:    `{"email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required email",
		},
		{
			inputJSON:    `{"email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "required password",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required email",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequestWithContext(context.Background(), "POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := adaptor.FiberHandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.Bytes(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			fmt.Println(responseMap["data"].(map[string]interface{})["data"])
			assert.Equal(t, responseMap["data"].(map[string]interface{})["data"], v.errorMessage)
		}
	}
}
