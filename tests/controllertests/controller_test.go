package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/controllers"
	"github.com/PixDale/sh-code-challenge/api/models"
)

var (
	server             = controllers.Server{}
	userInstance       = models.User{}
	taskInstance       = models.Task{}
	managerTokenJWT    = ""
	technicianTokenJWT = ""
	managerUser        = models.User{}
	technicianUser     = models.User{}
)

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()
	Authenticate()

	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDBDriver := "mysql"

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
	server.DB, err = gorm.Open(TestDBDriver, DBURL)
	if err != nil {
		fmt.Println(DBURL)
		fmt.Printf("Cannot connect to %s database: %s\n", TestDBDriver, err.Error())
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("[ControllerTest] We are connected to the %s database\n", TestDBDriver)
	}
}

func Authenticate() {
	err := refreshUserTable()
	if err != nil {
		log.Fatal("Failed to refresh user table:", err.Error())
	}
	userPass := "123"
	managerUser = models.User{
		Name:      "Felipe Galdino",
		Email:     "felipegaldino16@gmail.com",
		Password:  userPass,
		Role:      auth.ManagerRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	managerUser.Prepare()
	_, err = managerUser.SaveUser(server.DB)
	if err != nil {
		log.Fatal("Failed to create manager user:", err.Error())
	}

	managerTokenJWT, err = server.SignIn(managerUser.Email, userPass)
	if err != nil {
		log.Fatal("Failed to obtain an authentication token:", err.Error())
	}

	technicianUser = models.User{
		Name:      "PixDale",
		Email:     "pixdale@pixdale.com",
		Password:  userPass,
		Role:      auth.TechnicianRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	technicianUser.Prepare()
	_, err = technicianUser.SaveUser(server.DB)
	if err != nil {
		log.Fatal("Failed to create technician user:", err.Error())
	}

	technicianTokenJWT, err = server.SignIn(technicianUser.Email, userPass)
	if err != nil {
		log.Fatal("Failed to obtain an authentication token:", err.Error())
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
	log.Println("Successfully refreshed user table")
	return nil
}

func seedOneUser() (models.User, error) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Name:     "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}
	user.Prepare()
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {
	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		{
			Name:     "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Name:     "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}
	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
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
		Name:     "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
		Role:     auth.ManagerRole,
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Task{}, err
	}
	task := models.Task{
		Summary: "This is the content sam",
		UserID:  user.ID,
	}
	task.EncryptSummary()
	err = server.DB.Model(&models.Task{}).Create(&task).Error
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func seedUsersAndTasks() ([]models.User, []models.Task, error) {
	var err error

	users := []models.User{
		{
			Name:     "David Cossette",
			Email:    "david.cossette@gmail.com",
			Password: "password",
			Role:     auth.ManagerRole,
		},
		{
			Name:     "Mary Robbins",
			Email:    "mary.robbins@gmail.com",
			Password: "password",
			Role:     auth.TechnicianRole,
		},
	}
	tasks := []models.Task{
		{
			Summary: "Hello world 1",
		},
		{
			Summary: "Hello world 2",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tasks[i].UserID = users[i].ID
		tasks[i].EncryptSummary()

		err = server.DB.Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tasks table: %v", err)
		}
	}
	return users, tasks, nil
}
