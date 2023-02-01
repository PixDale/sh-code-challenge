package controllers

import (
	"github.com/PixDale/sh-code-challenge/api/middlewares"
)

func (server *Server) initializeRoutes() {
	// Home Route
	server.Router.Get("/", middlewares.SetMiddlewareJSON, server.Home)

	// Login Route
	server.Router.Post("/login", middlewares.SetMiddlewareJSON, server.Login)

	// Seed Route
	server.Router.Get("/seed", server.Seed)

	// Users routes
	server.Router.Post("/users", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.CreateUser)
	server.Router.Get("/users", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.GetUsers)
	server.Router.Get("/users/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.GetUser)
	server.Router.Put("/users/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.UpdateUser)
	server.Router.Delete("/users/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.DeleteUser)

	// Tasks routes
	server.Router.Post("/tasks", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.CreateTask)
	server.Router.Get("/tasks", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.GetTasks)
	server.Router.Get("/tasks/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.GetTask)
	server.Router.Put("/tasks/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.UpdateTask)
	server.Router.Delete("/tasks/:id", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, server.DeleteTask)
}
