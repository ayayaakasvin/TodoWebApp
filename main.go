package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"todoweb/packages/handlers"
	"todoweb/packages/handlers/authentication"
	"todoweb/packages/handlers/middleware"
	"todoweb/packages/handlers/task"
	"todoweb/packages/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	router   *gin.Engine
	database *utils.DataBaseProps
	store *sessions.CookieStore
	host string
	port string
)

func init() {
	router = gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*.html")

	store = sessions.NewCookieStore(utils.GenerateRandomKey(32))
	store.Options = &sessions.Options{
		Path: "/user",
		MaxAge: 300,
		HttpOnly: true,
		Secure: (os.Getenv("ENV") == "production"),
		SameSite: http.SameSiteLaxMode,
	}
	gob.Register(&utils.User{})

	err := godotenv.Load("database.env", "host.env")
    if err != nil {
        log.Fatal("Error loading database.env file")
    }

	dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

	host = os.Getenv("HOST")
	port = os.Getenv("PORT")

	database, err = utils.NewDatabaseConnection(dbName, dbHost, dbPort, dbUser, dbPassword)
	if err != nil || database == nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
}

func main() {
	AuthenticationHandlers := authentication.NewAuthenticationHandler(database, store)
	TaskHandlers := task.NewTaskHandler(database, store)
	MiddlewareHandlers := middleware.NewMiddlewareHandler(store)

	router.GET(handlers.RoutesPointer.MainLoginConfig.EmptyPathString, AuthenticationHandlers.GetEmptyPath)
	router.GET("/login", AuthenticationHandlers.GetLogin)
	router.POST("/login", AuthenticationHandlers.PostLogin)
	router.GET("/register", AuthenticationHandlers.GetRegister)
	router.POST("/register", AuthenticationHandlers.PostRegister)

	userRoutes := router.Group("/user", MiddlewareHandlers.Auth)
	{
		userRoutes.GET("/tasks", TaskHandlers.GetTasks)
		userRoutes.POST("/addTask", TaskHandlers.CreateTask)
		userRoutes.POST("/deleteTask", TaskHandlers.DeleteTask)
		userRoutes.POST("/logout", MiddlewareHandlers.Logout)
	}

	err := router.Run(host + ":" + port)
	if err != nil {
		log.Printf("Server running on %s:%s\n Error : %v", host, port, err)
	}
}

/*
	This project is huge victory over my laziness and some thoughts about suicide. However I am not good programmist,
	I am trying to improve my coding skills. Good luck to everyone !! :)
*/