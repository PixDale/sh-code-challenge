// Package responses provides structured HTTP responses for a Fiber application
package responses

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// UserResponse represents the response for a request
type UserResponse struct {
	// Status is the HTTP status code
	Status int `json:"status"`
	// Message is the response message
	Message string `json:"message"`
	// Data is the response data
	Data *fiber.Map `json:"data"`
}

// Error messages used in the application
var (
	ErrorFailReadToken     = errors.New("failed to read authorization token")
	ErrorUserNotAuthorized = errors.New("user not authorized")
	ErrorFailReadUserID    = errors.New("failed to read user id from token")
	ErrorTaskNotFound      = errors.New("task not found")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorFailParseClaims   = errors.New("failed to parse token claims")
)
