// Pakcage controllers implements several handler for the RestAPI app
package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/models"
	"github.com/PixDale/sh-code-challenge/api/responses"
	"github.com/PixDale/sh-code-challenge/api/utils/formaterror"
)

// CreateUser is a handler for creating a new user.
func (server *Server) CreateUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	user := models.User{}
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if !auth.HasRoleManager(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	user.Prepare()
	if validationErr := user.Validate(""); validationErr != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError.Error()}})
	}
	return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "success", Data: &fiber.Map{"data": userCreated}})
}

// GetUsers is a handler for retrieving a list of users.
func (server *Server) GetUsers(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	user := models.User{}
	defer cancel()

	if !auth.HasRoleManager(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

// GetUser is a handler for retrieving a specific user.
func (server *Server) GetUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	userID := c.Params("id")
	defer cancel()

	if !auth.HasRoleManager(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": userGotten}})
}

// UpdateUser is a handler for updating a user.
func (server *Server) UpdateUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	userID := c.Params("id")
	defer cancel()

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var tokenID uint32
	switch {
	case auth.HasRoleManager(c):
		break
	case auth.HasRoleTechnician(c):
		tokenID, err = auth.ExtractTokenID(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		if tokenID != uint32(uid) {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
		}
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError.Error()}})
	}
	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

// DeleteUser is a handler for deleting a user
func (server *Server) DeleteUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	userID := c.Params("id")
	defer cancel()

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	var tokenID uint32
	switch {
	case auth.HasRoleManager(c):
		break
	case auth.HasRoleTechnician(c):
		tokenID, err = auth.ExtractTokenID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
		}
		if tokenID != uint32(uid) {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
		}
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	user := models.User{}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusNoContent).JSON(
		responses.UserResponse{Status: fiber.StatusNoContent, Message: "success", Data: &fiber.Map{"data": fmt.Sprintf("User %d successfully deleted!", uid)}},
	)
}
