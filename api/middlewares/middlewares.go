package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/notification"
	"github.com/PixDale/sh-code-challenge/api/responses"
)

// SetMiddlewareJSON defines a middleware to set the content type of a request to JSON
func SetMiddlewareJSON(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	return c.Next()
}

// SetMiddlewareAuthentication defines a middleware to verify the integrity of the token within the request
func SetMiddlewareAuthentication(c *fiber.Ctx) error {
	err := auth.TokenValid(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}
	return c.Next()
}

func SetMiddlewareSendNotification(c *fiber.Ctx) error {
	if auth.HasRoleTechnician(c) {
		var err error
		uid, err := auth.ExtractTokenID(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": errors.New("failed to send a notification").Error()},
			})
		}
		method := c.Method()
		url := c.BaseURL()

		message := fmt.Sprintf("User %d tried to perform a %s request in the URL: %s", uid, method, url)
		err = notification.PublishNotification([]byte(message))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
	}
	return c.Next()
}
