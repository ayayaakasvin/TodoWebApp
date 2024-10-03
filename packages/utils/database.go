package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Database table and columns name
const (
	tasksTableName = "tasks"
	tasksDescription = "description"
	tasksIsCompleted = "is_completed"
	tasksUserID = "user_id"
	tasksID = "id"
)


// RegisterForm represents html form POST struct for checking and adding to database
type RegisterForm struct {
	Username string
	Password string
	Re_Password string
}

// DataBaseProps holds the properties needed to connect and manipulate the database
type DataBaseProps struct {
	databaseName string
	host         string
	port         string
	user         string
	password     string
	Connection   *sql.DB
}

// Define a struct to return user ID and authentication status
type AuthResult struct {
	UserItself User
	IsAuthenticated bool
}

type Task struct {
	Description string
	TaskID string
}

type ToDoPassStruct struct {
	Tasks  []Task
	UserID int
}

// NewRegisterForm return FormStruct struct
func NewRegisterForm (username, password, re_password string) *RegisterForm {
	return &RegisterForm{
		Username: username,
		Password: password,
		Re_Password: re_password,
	}
}

// NewDatabaseConnection creates a new database connection and returns DataBaseProps and an error if there is an issue
func NewDatabaseConnection(databaseName, host, port, user, password string) (*DataBaseProps, error) {
	databaseProps := &DataBaseProps{
		databaseName: databaseName,
		host:         host,
		port:         port,
		user:         user,
		password:     password,
	}

	err := databaseProps.ConnectTodatabase()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Database connection is successful\n")

	return databaseProps, nil
}

// connString returns the connection string for the database
func (database *DataBaseProps) connString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s databasename=%s sslmode=disable", database.host, database.port, database.user, database.password, database.databaseName)
}

// ConnectTodatabase tries to connect to the database and sets the Connection field, returning any connection error
func (database *DataBaseProps) ConnectTodatabase() error {
	var err error
	database.Connection, err = sql.Open("postgres", database.connString())
	if err != nil {
		log.Printf("Error during connection: %v\n", err)
		return err
	}

	if err = database.Connection.Ping(); err != nil {
		log.Printf("Error during pinging database: %v", err)
		return err
	}
	
	return nil
}

// ExecuteScript executes a script (INSERT INTO, DROP COLUMN, etc.) and returns the number of affected rows, the last inserted ID, and any error
func (database *DataBaseProps) ExecuteScript(script string, args ...any) (int64, error) {
	result, err := database.Connection.Exec(script, args...)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rowsAffected, nil
}

// CheckUserNameAndPassword checks if a username and hashed password combination exists in the database
func (database *DataBaseProps) CheckUserNameAndPassword(Username string, Password string) (AuthResult, error) {
	if database == nil || database.Connection == nil {
		return AuthResult{}, fmt.Errorf("database connection is nil")
	}

	var (
		storedHashPassword string
		userID             int
	)

	// Update the query to select both the hashed password and user ID
	scriptToFindUser := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = $1 LIMIT 1", usersPasswordHashColumn, usersIDColumn, tableUsersNaming, usersUsernameColumn)
	row := database.Connection.QueryRow(scriptToFindUser, Username)
	
	err := row.Scan(&storedHashPassword, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return AuthResult{}, fmt.Errorf("not found") // No rows means user not found
		}
		return AuthResult{}, err
	}

	if ComparePassword(Password, storedHashPassword) {
		return AuthResult{IsAuthenticated: true}, nil
	}

	return AuthResult{IsAuthenticated: false}, nil
}

func (database *DataBaseProps) FetchUserByUsername (Username string) (User, error) {
	if database == nil || database.Connection == nil {
		return User{}, fmt.Errorf("database connection is nil")
	}

	var (
		user User = User{}
		createTime time.Time
	)

	scriptToFindUser := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s = $1 LIMIT 1",
		tableUsersNaming,
		usersUsernameColumn,
	)

	err := database.Connection.QueryRow(scriptToFindUser, Username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&createTime,
	)
	if err != nil {
        if err == sql.ErrNoRows {
            return User{}, fmt.Errorf(UserNotFound, Username)
        }
        return User{}, fmt.Errorf("failed to scan user: %v", err)
    }

	user.creationTime = createTime.Format("2006-01-02 15:04:05")
	
	return user, nil
}

func (database *DataBaseProps) DoesUserExist (Username string) (bool, error) {
	if database == nil || database.Connection == nil {
		return false, fmt.Errorf("database connection is nil")
	}

	var scriptToFindUser string = fmt.Sprintf("SELECT 1 FROM %s WHERE %s = $1 LIMIT 1", tableUsersNaming, usersUsernameColumn)
	row := database.Connection.QueryRow(scriptToFindUser, Username)
	
	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists == 1, nil
}

func (database *DataBaseProps) IsValidRegister (UserInput RegisterForm, data gin.H) error {
	exists, err := database.DoesUserExist(UserInput.Username)
	if err != nil {
		return fmt.Errorf(InternalErrorString)
	}

	if exists {
		data[ErrorUsernameHTML] = UserAlreadyExistError
		data[Form] = UserInput
		return fmt.Errorf(UserAlreadyExistError)
	}
	
	if err := IsValidUsername(UserInput.Username); err != nil {
		data[ErrorUsernameHTML] = UsernameContainsSpace
		data[Form] = UserInput
		return err
	}

	if err := IsValidPassword(UserInput.Password, UserInput.Re_Password); err != nil {
		data[ErrorPasswordHTML] = err.Error()
		data[Form] = UserInput
		return err
	}

	return nil
}

// CreateNewUser creates a new user in the database with the given username and password
func (database *DataBaseProps) CreateNewUser (Username, Password string) error {
	if database == nil || database.Connection == nil {
		return fmt.Errorf("database connection is nil")
	}

	hashedPassword, err := HashPassword(Password)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)", tableUsersNaming, usersUsernameColumn, usersPasswordHashColumn)
	_, err = database.Connection.Exec(query, Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

// Used for getting username || Used before implementing Cookies, can be used for debuggin
func (database *DataBaseProps) GetUsername (userID string) (string, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1", usersUsernameColumn, tableUsersNaming, usersIDColumn)

	var result string
	err := database.Connection.QueryRow(query, userID).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user found with ID %s", userID)
		}
		return "", fmt.Errorf("failed to retrieve username: %v", err)
	}

	return result, nil
}

// AddTask adds Task to database by userID
func (database *DataBaseProps) AddTask (userID string, task string) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES ($1, $2)`, tasksTableName, tasksUserID, tasksDescription)

	_, err := database.Connection.Exec(query, userID, task)
	if err != nil {
		return fmt.Errorf("insert into error : %v", err)
	}

	return nil
}

// Used for fetching Tasks of User by userID from database 
func (database *DataBaseProps) GetTasksFromDatabase (userID string) ([]Task, error) {
	if database == nil || database.Connection == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var query string = fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = $1", tasksDescription, tasksID, tasksTableName, tasksUserID)
	rows, err := database.Connection.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("row query error : %v", err)
	}
	defer rows.Close()

	var result []Task = []Task{}

	for rows.Next() {
		var description string
		var id string
		if err := rows.Scan(&description, &id); err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		result = append(result, Task{TaskID: id, Description: description})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return result, nil
}

// DeleteTask deletes task by id
func (database *DataBaseProps) DeleteTask(userID string, taskID int) error {
	if database == nil || database.Connection == nil {
		return fmt.Errorf("database connection is nil")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1 AND %s = $2", tasksTableName, tasksUserID, tasksID)
	_, err := database.Connection.Exec(query, userID, taskID)
	if err != nil {
		return fmt.Errorf("row delete error: %v", err)
	}

	return nil
}