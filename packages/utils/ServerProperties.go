package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Some const's being used for different tasks, should be comment next time :)
const (
    DigitsRegexp = "[0-9]"
    LatinRegexp = `[a-zA-Z]`
    MinLength = 8
    MaxLength = 64
    InternalErrorString = "Internal error occurred during validation"
    LatinError = "Password should contain at least one latin letter"
    DigitError = "Password should contain at least one digit"
    SpaceError = "Password should not contain space"
    PasswordNotMatch = "Passwords do not match"
    UserAlreadyExistError = "User already exists"
    UsernameContainsSpace = "Username should not contain space"
    ErrorPasswordHTML = "PasswordError"
    ErrorUsernameHTML = "UserExistError"
    ErrorLoginHTML = "LoginError"
    LoginError = "Incorrect username or password"
    Form = "Form"
    ConvertError = "StrToInt error"
    GetTaskError = "Task parsing error"
    UserNotFound = "User %s not found"
)

// Checks if gained password valid, if it is unvalid return error
func IsValidPassword (pword, repword string) (error) {
    passwordByted := []byte(pword)

    if pword != repword {
        return fmt.Errorf(PasswordNotMatch)
    }

    if len(passwordByted) < MinLength {
        return fmt.Errorf("Password length must be at least %d", MinLength)
    } else if len(passwordByted) > MaxLength {
        return fmt.Errorf("Password length must be at most %d", MaxLength)
    }

    if matched, err := regexp.Match(LatinRegexp, passwordByted); err != nil {
        return fmt.Errorf(InternalErrorString)
    } else if !matched {
        return fmt.Errorf(LatinError)
    }

    if digitValid, err := regexp.Match(DigitsRegexp, passwordByted); err != nil {
        return fmt.Errorf(InternalErrorString)
    } else if !digitValid {
        return fmt.Errorf(DigitError)
    }

    if ContainSpace(pword) {
        return fmt.Errorf(SpaceError)
    }

    return nil
}

// Checks if gained username valid, if it is unvalid return error
func IsValidUsername (uname string) (error) {
    if ContainSpace(uname) {
        return fmt.Errorf(UsernameContainsSpace)
    }

    return nil
}

// Returns true if string contains space(\s)
func ContainSpace (input string) bool {
    return strings.Contains(input, " ")
}

// Returns string with trimed beginning and end 
func TrimSpace (input string) string {
    return strings.TrimSpace(input)
}

// Converts string type to int, if fails, returns -1
func StrToInt (input string) (int) {
    result, err := strconv.Atoi(input)
    if err != nil {
        return -1
    }

    return result
}