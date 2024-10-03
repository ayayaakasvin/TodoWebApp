package handlers

import (
	"todoweb/packages/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"todoweb/packages/config"
)
// RoutesPointer holds the route configuration for the application.
var RoutesPointer *config.RouteConfig = config.Routes

// getUserFromSession retrieves the user information from the session.
func GetUserFromSession(c *gin.Context, Store *sessions.CookieStore) (*utils.User, bool) {
	session, err := Store.Get(c.Request, RoutesPointer.Cookie.Naming)
	if err != nil {
		return nil, false // If session retrieval fails, return nil and false.
	}

	userInterface, ok := session.Values[RoutesPointer.Cookie.UserInfoKey].(*utils.User)
	return userInterface, ok // Return the user info and whether it was successful.
}

// GetSessionAndUser retrieves the session and user information from the provided request context.
// It first attempts to retrieve the session using the provided CookieStore and session naming convention.
// If the session retrieval fails, it returns nil for both the session and user, and false to indicate failure.
// If successful, it extracts the user information from the session values using the UserInfoKey.
// The function returns the session, the user information (if available), and a boolean indicating success.
func GetSessionAndUser (c *gin.Context, Store *sessions.CookieStore) (*sessions.Session, *utils.User, bool) {
    // Retrieve the session from the request using the session name defined in the Routes configuration.
    session, err := Store.Get(c.Request, RoutesPointer.Cookie.Naming)
    if err != nil {
        // If session retrieval fails, return nil for both the session and user, and false to indicate an error.
        return nil, nil, false
    }

    // Attempt to retrieve the user information from the session values using the predefined key.
    userInterface, ok := session.Values[RoutesPointer.Cookie.UserInfoKey].(*utils.User)
    
    // Return the session object, the user information, and a boolean indicating whether the retrieval was successful.
    return session, userInterface, ok
}

// GetSession retrieves the session associated with the current request.
// Parameters:
// - c: The Gin context, which contains the HTTP request and response.
// - Store: The session cookie store used to retrieve and manage session data.
// Returns:
// - *sessions.Session: The session object if successfully retrieved.
// - bool: A boolean indicating whether the session retrieval was successful (true) or if an error occurred (false).
func GetSession(c *gin.Context, Store *sessions.CookieStore) (*sessions.Session, bool) {
	// Retrieve the session from the request using the session name defined in the Routes configuration.
	session, err := Store.Get(c.Request, RoutesPointer.Cookie.Naming)
	if err != nil {
		// If session retrieval fails, return nil for the session and false to indicate an error.
		return nil, false
	}

	// Return the session object and a boolean indicating whether the retrieval was successful.
	return session, true
}

// SetSession sets a key-value pair in the session and saves it.
// Parameters:
// - c: The context from Gin, containing the HTTP request and response.
// - Store: The session cookie store used to manage session data.
// - key: The key for storing the session value (e.g., user ID, user info).
// - value: The value to be associated with the key in the session (e.g., user struct, user ID).
// Returns:
// - error: If an error occurs while retrieving or saving the session, it's returned, otherwise nil.
func SetSession(c *gin.Context, Store *sessions.CookieStore, key string, value interface{}) error {
	// Retrieve the session using the session name from the Routes configuration.
	session, err := Store.Get(c.Request, RoutesPointer.Cookie.Naming)
	if err != nil {
		// Return an error if there was an issue retrieving the session.
		return err
	}

	// Set the session value using the provided key and value.
	session.Values[key] = value

	// Save the session to persist the changes.
	err = sessions.Save(c.Request, c.Writer)
	if err != nil {
		// Return an error if saving the session fails.
		return err
	}

	// Return nil if no errors occurred.
	return nil
}
