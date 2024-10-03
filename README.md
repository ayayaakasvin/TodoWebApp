# TodoWebApp

TodoWebApp is a web application built using Go and the Gin framework that allows users to manage their to-do lists. This application provides a user-friendly interface for creating, viewing, and deleting tasks while ensuring secure user authentication.

## Features

- **User Authentication:** 
  - Users can register and log in securely.
  - Session management through cookies.

- **Task Management:**
  - Create new tasks.
  - View existing tasks.
  - Delete tasks as needed.

- **Dynamic HTML Rendering:**
  - Uses HTML templates to render pages for login, registration, and task management.

- **Middleware Integration:**
  - Protects user-specific routes to ensure only authenticated users can access their tasks.

- **Environment Configuration:**
  - Easily configurable through environment variables for deployment across different environments.

## Architecture

### Main Packages

- **Main Package (`main.go`):**
  - Initializes the router, middleware, and handlers.
  - Loads environment variables for database connection.

- **Handlers:**
  - **Authentication Handlers:** Manages user login, registration, and session handling.
  - **Task Handlers:** Handles operations related to tasks, such as fetching, creating, and deleting tasks.
  - **Middleware Handlers:** Implements authentication checks and other middleware functionalities.

- **Utilities:**
  - Helper functions and types for database connection management, user definitions, and session handling.

## Getting Started

### Prerequisites

- Go (version 1.23.0)
- Gin framework
- Gorilla sessions
- Godotenv
- PostgreSQL or another supported database

### Installation

1. Clone the repository:
    ```bash
    git clone <https://github.com/ayayaakasvin/TodoWebApp>
    cd TodoWebApp
    ```
2. Run the tidy command:
    ```bash
    go mod tidy
    ```
3. Run the `main.go` file:
    ```bash
    go run ./main.go
    ```

## Database

The application uses PostgreSQL as the database. If you want to change any database connection fields, you can edit the `database.env` file.

### "users" Table Structure

| Column Name    | Type       | Constraints                                   |
|----------------|------------|-----------------------------------------------|
| id             | int        | primary key, not null, default: auto-increment via `nextval('users_id_seq')` |
| username       | string     | not null                                     |
| passwordhash   | string     | not null                                     |
| creation_time  | time.Time  | default: current time via `now()`            |

### "tasks" Table Structure

| Column Name    | Type       | Constraints                                   |
|----------------|------------|-----------------------------------------------|
| id             | integer    | Primary Key, Not NULL, Default: `nextval('tasks_id_seq'::regclass)` |
| user_id        | integer    | Not NULL                                     |
| description    | character varying | length 255, Not NULL                  |
| is_completed    | boolean    | Default: false                               |
| created_at     | timestamp without time zone | Not NULL, Default: `CURRENT_TIMESTAMP` |

## Contact

You can contact me to: [kozhamseitov06@gmail.com](mailto:kozhamseitov06@gmail.com).

## End

Thank you for checking out TodoWebApp!