package controllers

import (
	"fmt"
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"

	"github.com/PixDale/sh-code-challenge/api/models"
)

const requestTimeout = 10 * time.Second

type Server struct {
	DB     *gorm.DB
	Router *fiber.App
}

func (server *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	var Dbdriver string = "mysql"

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Task{}) // database migration

	server.Router = fiber.New()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(server.Router.Listen(addr))
}
