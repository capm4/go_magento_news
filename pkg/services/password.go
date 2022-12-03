package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePassword(password string) error {
	passwordToCheck := strings.TrimSpace(password)
	if len(passwordToCheck) < 5 {
		return errors.New("password min length 5 char")
	}
	return nil
}
