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

// CreateTask implements the handler for the POST method of /tasks endpoint
func (server *Server) CreateTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	task := models.Task{}
	defer cancel()

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	task.Prepare()
	if validationErr := task.Validate(); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorFailReadToken.Error()}})
	}
	if uid != task.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}
	taskCreated, err := task.SaveTask(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError.Error()}})
	}
	taskCreated.DecryptSummary()
	return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "success", Data: &fiber.Map{"data": taskCreated}})
}

// GetTasks implements the handler for the GET method of /tasks endpoint.
// In case the given token contains the ManagerRole, it returns a list with all tasks.
// In case the given token contains the TechnicianRole, it returns a list of tasks created by this user
func (server *Server) GetTasks(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	task := models.Task{}
	defer cancel()
	var tasks *[]models.Task
	var err error

	switch {
	case auth.HasRoleManager(c):
		tasks, err = task.FindAllTasks(server.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	case auth.HasRoleTechnician(c):
		uid, err := auth.ExtractTokenID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorFailReadUserID.Error()}})
		}
		tasks, err = task.FindAllTasksByUserID(server.DB, uid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	// Decrypt summaries of tasks
	for _, t := range *tasks {
		t.DecryptSummary()
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": tasks}},
	)
}

func (server *Server) GetTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	taskID := c.Params("id")
	tid, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	task := models.Task{}
	taskReceived, err := task.FindTaskByID(server.DB, tid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	switch {
	case auth.HasRoleManager(c):
		break
	case auth.HasRoleTechnician(c):
		uid, err := auth.ExtractTokenID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorFailReadToken.Error()}})
		}
		if uid != taskReceived.UserID {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorUserNotAuthorized.Error()}})
		}
	}

	// Decrypt Summary of task
	taskReceived.DecryptSummary()
	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": taskReceived}})
}

func (server *Server) UpdateTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	taskID := c.Params("id")
	defer cancel()

	// Check if the task id is valid
	tid, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorFailReadToken.Error()}})
	}

	// Check if the task exist
	task := models.Task{}
	err = server.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&task).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.UserResponse{Status: fiber.StatusNotFound, Message: "error", Data: &fiber.Map{"data": responses.ErrorTaskNotFound.Error()}})
	}

	taskUpdate := models.Task{}
	// Read the task data
	if err = c.BodyParser(&taskUpdate); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	switch {
	case auth.HasRoleManager(c):
		break
	case auth.HasRoleTechnician(c):
		// If a user attempt to update a task not belonging to him
		if uid != task.UserID || uid != taskUpdate.UserID {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorUserNotAuthorized.Error()}})
		}
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorUserNotAuthorized.Error()}})
	}

	taskUpdate.Prepare()
	err = taskUpdate.Validate()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.UserResponse{Status: fiber.StatusUnprocessableEntity, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	taskUpdate.ID = task.ID // this is important to tell the model the task id to update, the other update field are set above

	taskUpdated, err := taskUpdate.UpdateATask(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": formattedError.Error()}})
	}
	taskUpdated.DecryptSummary()
	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": taskUpdated}})
}

func (server *Server) DeleteTask(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), requestTimeout)
	taskID := c.Params("id")
	defer cancel()

	// Is a valid task id given?
	tid, err := strconv.ParseUint(taskID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
	}

	// Check if the task exist
	task := models.Task{}
	err = server.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&task).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.UserResponse{Status: fiber.StatusNotFound, Message: "error", Data: &fiber.Map{"data": responses.ErrorTaskNotFound.Error()}})
	}

	switch {
	case auth.HasRoleManager(c):
		break
	case auth.HasRoleTechnician(c):
		// Is the authenticated user, the owner of this task?
		if uid != task.UserID {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": errors.New(http.StatusText(fiber.StatusUnauthorized)).Error()}})
		}

	default:
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": responses.ErrorUserNotAuthorized.Error()}})
	}

	_, err = task.DeleteATask(server.DB, tid, uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusNoContent).JSON(
		responses.UserResponse{Status: fiber.StatusNoContent, Message: "success", Data: &fiber.Map{"data": fmt.Sprintf("Task %d successfully deleted!", tid)}},
	)
}
