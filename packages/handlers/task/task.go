package task

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"net/http"
	"todoweb/packages/utils"
	"todoweb/packages/handlers"
)

/*
Tasks Table Structure:

- id (integer, Primary Key, Not NULL, Default: nextval('tasks_id_seq'::regclass))
  The unique identifier for each task.

- user_id (integer, Not NULL)
  The identifier for the user to whom the task belongs.

- description (character varying, length 255, Not NULL)
  A brief description of the task.

- is_completed (boolean, Default: false)
  Indicates whether the task has been completed (true) or not (false).

- created_at (timestamp without time zone, Not NULL, Default: CURRENT_TIMESTAMP)
  The timestamp when the task was created.
*/

// TaskHandlers interface defines the methods for task management.
type TaskHandlers interface {
	CreateTask(c *gin.Context) // Handles task creation.
	DeleteTask(c *gin.Context) // Handles task deletion.
	GetTasks(c *gin.Context)    // Retrieves tasks for the logged-in user.
}

// taskHandleProps struct holds dependencies for task handlers.
type taskHandleProps struct {
	Database *utils.DataBaseProps  // Database connection properties.
	Store    *sessions.CookieStore     // Cookie store for session management.
}

// GetTasks retrieves tasks for the authenticated user and renders the task page.
func (prop *taskHandleProps) GetTasks(c *gin.Context) {
	userInterface, ok := handlers.GetUserFromSession(c, prop.Store)
	if !ok {
		c.Redirect(http.StatusUnauthorized, handlers.RoutesPointer.UserConfig.GetTask.RedirectPath)
		return // Redirect to login if user is not authenticated.
	}

	UserTasks, err := prop.Database.GetTasksFromDatabase(userInterface.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return // Handle error if task retrieval fails.
	}

	c.HTML(http.StatusOK, handlers.RoutesPointer.UserConfig.GetTask.HTMLPageName, gin.H{
		"tasks": utils.ToDoPassStruct{
			Tasks:  UserTasks,                       // Pass the retrieved tasks to the template.
			UserID: utils.StrToInt(userInterface.ID), // Pass user ID for reference.
		},
		"Username": userInterface.Username, // Pass the username for display.
	})
}

// CreateTask handles the creation of a new task.
func (prop *taskHandleProps) CreateTask(c *gin.Context) {
	userInterface, ok := handlers.GetUserFromSession(c, prop.Store)
	if !ok {
		c.Redirect(http.StatusUnauthorized, handlers.RoutesPointer.UserConfig.GetTask.RedirectPath)
		return // Redirect to login if user is not authenticated.
	}

	if err := c.Request.ParseForm(); err != nil {
		c.Redirect(http.StatusSeeOther, handlers.RoutesPointer.UserConfig.GetTask.Route)
		return // Handle error if form parsing fails.
	}

	task := utils.TrimSpace(c.PostForm("taskTitle")) // Get and trim the task title.

	// Add the task to the database and handle any errors.
	if err := prop.Database.AddTask(userInterface.ID, task); err != nil {
		c.String(http.StatusInternalServerError, "Failed to add task")
		return // Handle error if task addition fails.
	}

	c.Redirect(http.StatusFound, handlers.RoutesPointer.UserConfig.DeleteTask.RedirectPath) // Redirect after successful creation.
}

// DeleteTask handles the deletion of a task.
func (prop *taskHandleProps) DeleteTask(c *gin.Context) {
	userInterface, ok := handlers.GetUserFromSession(c, prop.Store)
	if !ok {
		c.Redirect(http.StatusUnauthorized, handlers.RoutesPointer.UserConfig.GetTask.RedirectPath)
		return // Redirect to login if user is not authenticated.
	}

	if err := c.Request.ParseForm(); err != nil {
		c.Redirect(http.StatusSeeOther, handlers.RoutesPointer.UserConfig.GetTask.RedirectPath)
		return // Handle error if form parsing fails.
	}

	taskIDstr := utils.TrimSpace(c.PostForm("TaskID")) // Get and trim the task ID.
	taskID := utils.StrToInt(taskIDstr) // Convert task ID string to integer.
	if taskID == -1 {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return // Handle error if task ID conversion fails.
	}

	// Attempt to delete the task from the database and handle any errors.
	if err := prop.Database.DeleteTask(userInterface.ID, taskID); err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete task")
		return // Handle error if task deletion fails.
	}

	c.Redirect(http.StatusFound, handlers.RoutesPointer.UserConfig.DeleteTask.RedirectPath) // Redirect after successful deletion.
}

// NewTaskHandler creates a new instance of TaskHandlers with the provided database and session store.
func NewTaskHandler(db *utils.DataBaseProps, store *sessions.CookieStore) TaskHandlers {
	return &taskHandleProps{
		Database: db,  // Set the database property.
		Store:    store, // Set the session store property.
	}
}