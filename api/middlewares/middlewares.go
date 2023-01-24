package middlewares

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/PixDale/sh-code-challenge/api/auth"
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
