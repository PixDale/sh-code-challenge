// Package seed provides seed data for the database and functions to
// load the data into the database
package seed

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/PixDale/sh-code-challenge/api/auth"
	"github.com/PixDale/sh-code-challenge/api/models"
)

// users is an array of User structs containing seed data for the users table
var users = []models.User{
	{
		Name:     "Felipe Galdino",
		Email:    "felipegaldino16@gmail.com",
		Password: "123",
		Role:     auth.ManagerRole,
	},
	{
		Name:     "David Cossette",
		Email:    "david.cossette@gmail.com",
		Password: "123",
		Role:     auth.TechnicianRole,
	},
	{
		Name:     "Mary Robbins",
		Email:    "mary.robbins@gmail.com",
		Password: "123",
		Role:     auth.TechnicianRole,
	},
}

// tasks is an array of Task structs containing seed data for the tasks table
var tasks = []models.Task{
	{
		Summary: "Hello world 4",
	},
	{
		Summary: "Hello world 5",
	},
	{
		Summary: "Hello world 6",
	},
}

// Load is a function to load the seed data into the database
// It takes a pointer to a gorm.DB instance as an argument
func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Task{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Task{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Task{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tasks[i].UserID = users[i].ID
		tasks[i].EncryptSummary()

		err = db.Debug().Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tasks table: %v", err)
		}
	}
}
