package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/models"
	"github.com/PixDale/sh-code-challenge/api/responses"
	"github.com/PixDale/sh-code-challenge/api/utils/formaterror"
)

func (server *Server) CreateUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	user := models.User{}
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	user.Prepare()
	// use the validator library to validate required fields
	if validationErr := user.Validate(""); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": userCreated}})
}

func (server *Server) GetUsers(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	user := models.User{}
	defer cancel()

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

func (server *Server) GetUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	uid, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": userGotten}})
}

func (server *Server) UpdateUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	uid, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: http.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if tokenID != uint32(uid) {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(http.StatusUnauthorized))}})
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: http.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func (server *Server) DeleteUser(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	uid, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("unauthorized").Error()}})
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(http.StatusUnauthorized))}})
	}

	user := models.User{}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusNoContent).JSON(
		responses.UserResponse{Status: http.StatusNoContent, Message: "success", Data: &fiber.Map{"data": fmt.Sprintf("User %d successfully deleted!", uid)}},
	)
}
