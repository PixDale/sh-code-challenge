package controllers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/models"
	"github.com/PixDale/sh-code-challenge/api/responses"
	"github.com/PixDale/sh-code-challenge/api/utils/formaterror"
)

// Login handles user login by parsing the user data from request body, validating it, and then signing in.
// If the validation fails or sign in fails, it returns appropriate error status and message.
// If successful, it returns a success status with a token.
func (server *Server) Login(c *fiber.Ctx) error {
	user := models.User{}

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": formattedError.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": token}},
	)
}

// SignIn signs in the user by checking the provided email and password against the database.
// If the email does not match, it returns an error.
// If the password does not match, it returns an error.
// If successful, it returns a token.
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
