package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("password is incorrect")
)

// ValidatePassword compare password with hash 
func ValidatePassword(password string, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))

	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}
