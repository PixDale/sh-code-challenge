package modeltests

import (
	"log"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"

	"github.com/PixDale/sh-code-challenge/api/models"
)

func TestFindAllTasks(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table %v\n", err)
	}
	_, _, err = seedUsersAndTasks()
	if err != nil {
		log.Fatalf("Error seeding user and task  table %v\n", err)
	}
	tasks, err := taskInstance.FindAllTasks(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the tasks: %v\n", err)
		return
	}
	assert.Equal(t, len(*tasks), 2)
}

func TestSaveTask(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error user and task refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newTask := models.Task{
		ID:      1,
		Summary: "This is the content",
		UserID:  user.ID,
	}
	savedTask, err := newTask.SaveTask(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the task: %v\n", err)
		return
	}
	assert.Equal(t, newTask.ID, savedTask.ID)
	assert.Equal(t, newTask.Summary, savedTask.Summary)
	assert.Equal(t, newTask.UserID, savedTask.UserID)
}

func TestGetTaskByID(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table: %v\n", err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundTask, err := taskInstance.FindTaskByID(server.DB, task.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundTask.ID, task.ID)
	assert.Equal(t, foundTask.Summary, task.Summary)
}

func TestUpdateATask(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table: %v\n", err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	taskUpdate := models.Task{
		ID:      1,
		Summary: "modiupdate@gmail.com",
		UserID:  task.UserID,
	}
	updatedTask, err := taskUpdate.UpdateATask(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedTask.ID, taskUpdate.ID)
	assert.Equal(t, updatedTask.Summary, taskUpdate.Summary)
	assert.Equal(t, updatedTask.UserID, taskUpdate.UserID)
}

func TestDeleteATask(t *testing.T) {
	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table: %v\n", err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := taskInstance.DeleteATask(server.DB, task.ID, task.UserID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	// one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	// Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}

func TestEncryptAndDecryptSummary(t *testing.T) {
	task := models.Task{
		ID:        0,
		Summary:   "Task Decrypted",
		User:      models.User{},
		UserID:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	task.EncryptSummary()
	task.DecryptSummary()
}
