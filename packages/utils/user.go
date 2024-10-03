package utils

// Table: users
//
// Columns:
// 1. id (int, primary key, not null, default: auto-increment via nextval('users_id_seq'))
//    - This column is the primary key and auto-increments using a sequence in PostgreSQL.
//
// 2. username (string, not null)
//    - This column stores the username.
//
// 3. passwordhash (string, not null)
//    - This column stores the hashed password.
//
// 4. creation_time (time.Time, default: current time via now())
//    - This column stores the timestamp when the account was created, with a default value of the current time.

const (
	tableUsersNaming = "users"
	usersIDColumn = "id"
	usersUsernameColumn = "username"
	usersPasswordHashColumn = "passwordhash"
	usersCreationTimeColumn = "creation_time"
)

// Used for gob register and cookies, creationTime field is not used, in future I want to use them in tasks displaying
type User struct {
	Username string
	PasswordHash string
	ID string
	creationTime string
}