package utils

import "golang.org/x/crypto/bcrypt"
import "crypto/rand"
import "log"

// HashPassword hashes the given password using bcrypt and returns the hashed password as a string.
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    return string(hashedPassword), nil
}

// Compares password with hash from database
func ComparePassword (password, hash string) (bool) {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// Compares two hashes and returs true if they are the same
func CompareHash (hash1, hash2 string) (bool) {
	return hash1 == hash2
}

// GenerateRandomKey creates a random byte slice of the given size
func GenerateRandomKey(size int) []byte {
    key := make([]byte, size)
    if _, err := rand.Read(key); err != nil {
        log.Fatalf("Error generating random key: %v\n", err)
    }
    return key
}