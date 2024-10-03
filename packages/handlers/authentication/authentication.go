package authentication

import (
	"fmt"

	"github.com/gin-gonic/gin" // Gin framework for HTTP handling
	"github.com/gorilla/sessions" // Package for session management

	"net/http"

	"todoweb/packages/handlers" // Handlers for routing
	"todoweb/packages/utils" // Utility functions and types
)

// AuthenticationHandlers defines the interface for authentication-related handlers.
type AuthenticationHandlers interface {
	GetLogin(c *gin.Context) // Handler for GET login requests
	PostLogin(c *gin.Context) // Handler for POST login requests
	GetRegister(c *gin.Context) // Handler for GET registration requests
	PostRegister(c *gin.Context) // Handler for POST registration requests
	GetEmptyPath(c *gin.Context) // Handler for empty path redirects
}

// authenticationHandlerProps holds the properties needed for authentication handlers.
type authenticationHandlerProps struct {
	Database *utils.DataBaseProps  // Database connection properties.
	Store    *sessions.CookieStore  // Cookie store for session management.
}

// GetLogin renders the login page.
func (prop *authenticationHandlerProps) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, handlers.RoutesPointer.MainLoginConfig.PageName, nil)
}

// GetEmptyPath redirects to the main login page.
func (prop *authenticationHandlerProps) GetEmptyPath(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, handlers.RoutesPointer.MainLoginConfig.Path)
}

// PostLogin handles user login attempts.
func (prop *authenticationHandlerProps) PostLogin(c *gin.Context) {
	// Parse the form data from the request
	if err := c.Request.ParseForm(); err != nil {
		c.Redirect(http.StatusSeeOther, handlers.RoutesPointer.MainLoginConfig.EmptyPathString)
		return
	}

	// Retrieve and trim the username and password from the form
	username := utils.TrimSpace(c.PostForm(handlers.RoutesPointer.Authentication.ParseKeys.UsernameParseKey))
	password := utils.TrimSpace(c.PostForm(handlers.RoutesPointer.Authentication.ParseKeys.PasswordParseKey))

	// Fetch the user from the database using the username
	authResult, err := prop.Database.FetchUserByUsername(username)
	if err != nil {
		var errorMessage string = fmt.Sprintf(utils.UserNotFound, username)
		// Check if the error is due to the user not being found
		if err.Error() == errorMessage {
			data := gin.H{
				utils.ErrorLoginHTML: errorMessage, // Display the user not found error
				"Username":           username, // Pass the username back to the view
			}
			c.HTML(http.StatusOK, handlers.RoutesPointer.MainLoginConfig.PageName, data)
			return
		}

		// Handle other errors by returning an internal server error
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Check if the provided password matches the hashed password in the database
	if !utils.ComparePassword(password, authResult.PasswordHash) {
		data := gin.H{
			utils.ErrorLoginHTML: utils.LoginError, // Display login error
			"Username":           username, // Pass the username back to the view
		}
		c.HTML(http.StatusOK, handlers.RoutesPointer.MainLoginConfig.PageName, data)
		return
	}

	// Set the user session upon successful login
	handlers.SetSession(c, prop.Store, handlers.RoutesPointer.Cookie.UserInfoKey, authResult)

	// Redirect to the user's tasks page
	c.Redirect(http.StatusFound, handlers.RoutesPointer.UserConfig.GetTask.Route)
}

// GetRegister renders the registration page.
func (prop *authenticationHandlerProps) GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, handlers.RoutesPointer.MainRegisterConfig.PageName, nil)
}

// PostRegister handles user registration attempts.
func (prop *authenticationHandlerProps) PostRegister(c *gin.Context) {
	// Parse the form data from the request
	if err := c.Request.ParseForm(); err != nil {
		c.Redirect(http.StatusSeeOther, "/register")
		return
	}

	// Create a new register form with trimmed inputs
	registerForm := utils.NewRegisterForm(
		utils.TrimSpace(c.PostForm("uname")), 
		utils.TrimSpace(c.PostForm("pword")), 
		utils.TrimSpace(c.PostForm("re-pword")), 
	)

	var (
		data   gin.H = gin.H{} // Data to pass to the HTML template
		status int // Status code for response
	)

	// Validate the registration form
	if err := prop.Database.IsValidRegister(*registerForm, data); err != nil {
		c.HTML(status, handlers.RoutesPointer.MainRegisterConfig.PageName, data) // Render errors if validation fails
		return
	}

	// Create a new user in the database
	if err := prop.Database.CreateNewUser(registerForm.Username, registerForm.Password); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error") // Handle errors during user creation
	}

	// Redirect to the login page after successful registration
	c.Redirect(http.StatusSeeOther, handlers.RoutesPointer.MainLoginConfig.Path)
}

// NewAuthenticationHandler creates a new instance of AuthenticationHandlers.
func NewAuthenticationHandler(db *utils.DataBaseProps, store *sessions.CookieStore) AuthenticationHandlers {
	return &authenticationHandlerProps{
		Database: db,  // Set the database property.
		Store:    store, // Set the session store property.
	}
}