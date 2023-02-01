package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/PixDale/sh-code-challenge/api/responses"
	"github.com/PixDale/sh-code-challenge/api/seed"
)

// Seed handles the GET /seed request, to refresh the database with the seed values
func (server *Server) Seed(c *fiber.Ctx) error {
	seed.Load(server.DB)

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": "called seeder"}})
}
