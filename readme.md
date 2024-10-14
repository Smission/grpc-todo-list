# gRPC Todo List Application

This is a gRPC-based Todo List application built with Go. It supports both an interactive command-line interface (CLI) and a REST API (via gRPC-Gateway) for managing tasks. The application uses a MySQL database and can be configured via a YAML file.

## Features

- gRPC server for handling task-related operations (add, get, complete)
- REST API via gRPC-Gateway
- Interactive CLI for managing tasks
- MySQL database for task storage
- Configurable via a YAML configuration file
- Dockerized for easy setup

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [CLI Usage](#cli-usage)
  - [REST API Usage](#rest-api-usage)
- [Running with Docker](#running-with-docker)
- [Commands](#commands)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/grpc-todo-list.git
    cd grpc-todo-list
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Ensure MySQL is running locally or via Docker.

4. Build the Go application:

    ```bash
    go build -o todo-app
    ```

## Configuration

Configure the application using the `config.yaml` file in the root of the project. Example configuration:

```yaml
env: development
grpc_port: 50054
db:
  user: root
  name: todo_app
  host: localhost
  password:
  port: 3306
```

## Database Configuration
The db section contains the MySQL connection details. Ensure the database exists, and provide the credentials in the YAML file.

## Environment Variables
Alternatively, you can use environment variables:

DB_USER
DB_PASSWORD
DB_HOST
DB_PORT
DB_NAME
GRPC_PORT
HTTP_PORT
Usage
CLI Usage

## Once the application is running, you can interact via the CLI:
./todo-app

# The CLI provides the following options (which can be used to used to do different operations):

[1] Add a Task
[2] Get all Tasks
[3] Complete a Task
[4] Exit


# Swagger Rest API Definition

``` yaml
openapi: 3.0.1
info:
  title: Todo List API
  description: API for managing tasks in a Todo List application.
  version: 1.0.0
servers:
  - url: http://localhost:4000
paths:
  /task:
    post:
      summary: Add a new task
      operationId: addTask
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                title:
                  type: string
                description:
                  type: string
              required:
                - title
      responses:
        '200':
          description: Task added successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                    description: The ID of the created task
        '400':
          description: Invalid input
          
  /task/all:
    get:
      summary: Retrieve all tasks
      operationId: getTasks
      responses:
        '200':
          description: A list of tasks
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                    title:
                      type: string
                    description:
                      type: string
                    status:
                      type: string
                      enum: [Incomplete, Completed]
        '500':
          description: Internal server error

  /task/{id}:
    patch:
      summary: Mark a task as complete
      operationId: completeTask
      parameters:
        - name: id
          in: path
          required: true
          description: ID of the task to complete
          schema:
            type: integer
      responses:
        '200':
          description: Task completed successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    description: Indicates if the operation was successful
        '404':
          description: Task not found
        '400':
          description: Invalid ID format
```

