# Sword Health Code Challenge

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![go.mod](https://img.shields.io/github/go-mod/go-version/PixDale/sh-code-challenge/main)](go.mod)
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

**PS:** If needed you can clean docker cache before start, using the command **make docker_clean**


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
To run the test environment along with the unit tests run the command:
```
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

### **Doc**
To check the Golang documentation, first you need to have **godoc** cmd. To install it run:
```
go install -v golang.org/x/tools/cmd/godoc@latest
```

After installed, run:
```
godoc --http=localhost:6060
```
then access: **http://127.0.0.1:6060/pkg/github.com/PixDale/sh-code-challenge/**



### **Kubernetes**
To start the kubernetes deployment, first you need to have **`minikube`** and **`kubectl`** installed, then run:
```
make kube_start
make kube_apply
```

After that if you want to stop, run:
```
make kube_stop
```
## **License**

This project is licensed under the MIT License.