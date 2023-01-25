package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	"github.com/PixDale/sh-code-challenge/api/controllers"
	"github.com/PixDale/sh-code-challenge/api/models"
)

var (
	server       = controllers.Server{}
	userInstance = models.User{}
	taskInstance = models.Task{}
)

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDBDriver := "mysql"

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
	server.DB, err = gorm.Open(TestDBDriver, DBURL)
	if err != nil {
		fmt.Println(DBURL)
		fmt.Printf("Cannot connect to %s database\n", TestDBDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDBDriver)
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	refreshUserTable()

	user := models.User{
		Name:     "Felipe Galdino",
		Email:    "felipegaldino16@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {
	users := []models.User{
		{
			Name:     "David Cossette",
			Email:    "david.cossette@gmail.com",
			Password: "password",
		},
		{
			Name:     "Mary Robbins",
			Email:    "mary.robbins@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndTaskTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Task{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Task{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneTask() (models.Task, error) {
	err := refreshUserAndTaskTable()
	if err != nil {
		return models.Task{}, err
	}
	user := models.User{
		Name:     "Tom Holland",
		Email:    "th@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Task{}, err
	}
	task := models.Task{
		Summary: "Turned the lights off.",
		UserID:  user.ID,
	}
	err = server.DB.Model(&models.Task{}).Create(&task).Error
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func seedUsersAndTasks() ([]models.User, []models.Task, error) {
	var err error

	if err != nil {
		return []models.User{}, []models.Task{}, err
	}
	users := []models.User{
		{
			Name:     "David Cossette",
			Email:    "david.cossette@gmail.com",
			Password: "password",
		},
		{
			Name:     "Mary Robbins",
			Email:    "mary.robbins@gmail.com",
			Password: "password",
		},
	}
	tasks := []models.Task{
		{
			Summary: "Drunk water",
		},
		{
			Summary: "Slept",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tasks[i].UserID = users[i].ID

		err = server.DB.Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tasks table: %v", err)
		}
	}
	return users, tasks, nil
}
