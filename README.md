# Sword Health Code Challenge

[![Go Reference](https://pkg.go.dev/badge/github.com/PixDale/sh-code-challenge.svg)](https://pkg.go.dev/github.com/PixDale/sh-code-challenge)
[![go.mod](https://img.shields.io/github/go-mod/go-version/PixDale/sh-code-challenge)](go.mod)
[![LICENSE](https://img.shields.io/github/license/PixDale/sh-code-challenge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/PixDale/sh-code-challenge)](https://goreportcard.com/report/github.com/PixDale/sh-code-challenge)

This is a task management application written in Go that provides a REST API for authentication and authorization of users with different roles. The application and its dependencies can be run using Docker Compose.

## **Features**

- Task management: users can create, read, update, and delete tasks
- Role-based authentication: users are assigned different roles (manager and technician) and are authorized to perform certain actions based on their role
- REST API: the application provides a REST API for managing tasks, users and performing authentication and authorization
- Docker Compose: the application and its dependencies can be run using Docker Compose for easy setup and configuration

## **Prerequisites**

- Docker and Docker Compose

## **Getting Started**

### **Installation**

1. Clone the repository

```
git clone https://github.com/PixDale/sh-code-challenge.git
```

2. Navigate to the project directory

```
cd sh-code-challenge
```

3. Start the application using Docker Compose

```
docker-compose up
```

### **API Endpoints**

The API provides the following endpoints for task management and authentication:

- POST **`/tasks`**: create a new task
    - Requests must have a user with the Manager role
- GET **`/tasks`**: retrieve a list of tasks
    - For requests with the Manager role, all tasks will be retrieved
    - For requests with the Technician role, only the task from this user will be retrieved
- GET **`/tasks/{id}`**: retrieve a single task by ID
    - For requests with the Manager role, the task is retrieved unconditionally
    - For requests with the Technician role, the task is retrieved only if it belongs to this user
- PUT **`/tasks/{id}`**: update an existing task
    - For requests with the Manager role, the task is updated unconditionally
    - For requests with the Technician role, the task is updated only if it belongs to this user
- DELETE **`/tasks/{id}`**: delete a task
    - Requests must have a user with the Manager role

All user management requests must have a user with the Manager role

- POST **`/users`**: register a new user
- GET **`/users`**: retrieve a list of users
- GET **`/users/{id}`**: retrieve a single user by ID
- PUT **`/users/{id}`**: update an existing user
- DELETE **`/users/{id}`**: delete a user

- POST **`/login`**: login to retrieve a JWT token
    - Return the JWT Token containing information such as User ID and Role

### **Authentication**

Access to all endpoints except for login, requires a JSON Web Token (JWT) for authentication. The token must be included in the **`Authorization`** header of the request in the following format:

```
Authorization: Bearer [JWT_TOKEN]
```
### **Tests**
To run the tests...

## **License**

This project is licensed under the MIT License.
