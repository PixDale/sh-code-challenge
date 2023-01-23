package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/models"
	"github.com/PixDale/sh-code-challenge/api/responses"
	"github.com/PixDale/sh-code-challenge/api/utils/formaterror"
)

func (server *Server) Login(c *fiber.Ctx) error {
	user := models.User{}

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: http.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: http.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": formattedError}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": token}},
	)
}

func (server *Server) SignIn(email, password string) (string, error) {
	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID, user.Role)
}
