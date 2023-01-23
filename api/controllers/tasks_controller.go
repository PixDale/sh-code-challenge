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

func (server *Server) CreateTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	task := models.Task{}
	defer cancel()

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	task.Prepare()
	if validationErr := task.Validate(); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("failed to read authorization token").Error()}})
	}
	if uid != task.UserID {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(http.StatusUnauthorized))}})
	}
	taskCreated, err := task.SaveTask(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError}})
	}
	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": taskCreated}})
}

func (server *Server) GetTasks(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	task := models.Task{}
	defer cancel()

	tasks, err := task.FindAllTasks(server.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": tasks}},
	)
}

func (server *Server) GetTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskID := c.Params("taskId")
	defer cancel()

	tid, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	task := models.Task{}
	taskReceived, err := task.FindTaskByID(server.DB, tid)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": taskReceived}})
}

func (server *Server) UpdateTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskID := c.Params("taskId")
	defer cancel()

	// Check if the task id is valid
	tid, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("failed to read authorization token").Error()}})
	}

	// Check if the task exist
	task := models.Task{}
	err = server.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&task).Error
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": errors.New("task not found")}})
	}

	// If a user attempt to update a task not belonging to him
	if uid != task.UserID {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("user not authorized").Error()}})
	}

	taskUpdate := models.Task{}
	// Read the task data
	if err := c.BodyParser(&taskUpdate); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: http.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Also check if the request user id is equal to the one gotten from token
	// if uid != taskUpdate.UserID {
	// 	return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("unauthorized").Error()}})
	// }

	taskUpdate.Prepare()
	err = taskUpdate.Validate()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: http.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	taskUpdate.ID = task.ID // this is important to tell the model the task id to update, the other update field are set above

	taskUpdated, err := taskUpdate.UpdateATask(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": taskUpdated}})
}

func (server *Server) DeleteTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskID := c.Params("taskId")
	defer cancel()

	// Is a valid task id given?
	tid, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("unauthorized").Error()}})
	}

	// Check if the task exist
	task := models.Task{}
	err = server.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&task).Error
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": errors.New("Task not found")}})
	}

	// Is the authenticated user, the owner of this task?
	if uid != task.UserID {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New("unauthorized").Error()}})
	}
	_, err = task.DeleteATask(server.DB, tid, uid)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusNoContent).JSON(
		responses.UserResponse{Status: http.StatusNoContent, Message: "success", Data: &fiber.Map{"data": fmt.Sprintf("Task %d successfully deleted!", tid)}},
	)
}
