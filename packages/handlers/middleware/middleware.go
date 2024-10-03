package middleware

import (
	"github.com/gin-gonic/gin" 
	"github.com/gorilla/sessions"

	"net/http"
	"fmt"
	"todoweb/packages/handlers"
)

// MiddlewareHandlers defines the interface for authentication-related middleware.
// Auth: Ensures that users are authenticated.
// Logout: Logs out the user by terminating their session.
type MiddlewareHandlers interface {
	Auth(c *gin.Context)
	Logout(c *gin.Context)
}

// authHandler contains a CookieStore for session management.
type authHandler struct {
	Store *sessions.CookieStore
}

// Auth checks if the user is authenticated.
// If the user session is valid, it refreshes the session expiry time and proceeds to the next handler.
// If not authenticated, the user is redirected to the login page.
func (BrowserAuth *authHandler) Auth(c *gin.Context) {
	// Retrieve session and user information from the session store.
	session, _, ok := handlers.GetSessionAndUser(c, BrowserAuth.Store)
	if !ok {
		// If session or user info is missing, redirect to the login page.
		c.Redirect(http.StatusUnauthorized, handlers.RoutesPointer.UserConfig.GetTask.RedirectPath)
		return
	}

	// Refresh the session expiration time.
	session.Options.MaxAge = handlers.RoutesPointer.Authentication.SessionTime
	err := sessions.Save(c.Request, c.Writer)
	if err != nil {
		// Handle session saving error.
		c.String(http.StatusInternalServerError, err.Error())
		fmt.Printf("session save error: %v\n", err)
		return
	}
	
	// Call the next handler in the chain (if the user is authenticated).
	c.Next()
}

// Logout terminates the user session by setting its MaxAge to -1 (expire immediately).
// After successfully logging out, the user is redirected to the login page.
func (BrowserAuth *authHandler) Logout(c *gin.Context) {
	// Retrieve session and user information from the session store.
	session, _, ok := handlers.GetSessionAndUser(c, BrowserAuth.Store)
	if !ok {
		// If session or user info is missing, redirect to the login page.
		c.Redirect(http.StatusUnauthorized, handlers.RoutesPointer.UserConfig.GetTask.RedirectPath)
		return
	}

	// Invalidate the session by setting MaxAge to -1.
	session.Options.MaxAge = handlers.RoutesPointer.MainLoginConfig.SessionTimeOut
	err := sessions.Save(c.Request, c.Writer)
	if err != nil {
		// Handle session saving error.
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Redirect the user to the login page after logging out.
	c.Redirect(http.StatusFound, handlers.RoutesPointer.MainLoginConfig.Path) // Consider making this path configurable.
}

// NewMiddlewareHandler creates a new authHandler instance that implements the MiddlewareHandlers interface.
// It takes a CookieStore for session management.
func NewMiddlewareHandler(store *sessions.CookieStore) MiddlewareHandlers {
	return &authHandler{
		Store: store,
	}
}
