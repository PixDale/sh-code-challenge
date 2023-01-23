package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func (server *Server) Home(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Welcome To This API",
	})
}
