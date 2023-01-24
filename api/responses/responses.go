package responses

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

// Error messages
var (
	ErrorFailReadToken     = errors.New("failed to read authorization token")
	ErrorUserNotAuthorized = errors.New("user not authorized")
	ErrorFailReadUserID    = errors.New("failed to read user id from token")
	ErrorTaskNotFound      = errors.New("task not found")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorFailParseClaims   = errors.New("failed to parse token claims")
)
