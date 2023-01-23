package controllers

import "github.com/PixDale/sh-code-challenge/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.Get("/", middlewares.SetMiddlewareJSON, s.Home)

	// Login Route
	s.Router.Post("/login", middlewares.SetMiddlewareJSON, s.Login)

	// Users routes
	s.Router.Post("/users", middlewares.SetMiddlewareJSON, s.CreateUser)
	s.Router.Get("/users", middlewares.SetMiddlewareJSON, s.GetUsers)
	s.Router.Get("/users/:userId", middlewares.SetMiddlewareJSON, s.GetUser)
	s.Router.Put("/users/:userId", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, s.UpdateUser)
	s.Router.Delete("/users/:userId", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, s.DeleteUser)

	// Posts routes
	s.Router.Post("/tasks", middlewares.SetMiddlewareJSON, s.CreateTask)
	s.Router.Get("/tasks", middlewares.SetMiddlewareJSON, s.GetTasks)
	s.Router.Get("/tasks/:taskId", middlewares.SetMiddlewareJSON, s.GetTask)
	s.Router.Put("/tasks/:taskId", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, s.UpdateTask)
	s.Router.Delete("/tasks/:taskId", middlewares.SetMiddlewareJSON, middlewares.SetMiddlewareAuthentication, s.DeleteTask)
}
