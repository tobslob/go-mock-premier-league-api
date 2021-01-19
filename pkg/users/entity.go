package users

import (
	"time"
)

// User model
type User struct {
	ID          string    `json:"_id,omitempty" bson:"_id"`
	CreadAt     time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	FirstName   string    `json:"first_name,omitempty" bson:"first_name"`
	LastName    string    `json:"last_name,omitempty" bson:"last_name"`
	EmailAdress string    `json:"email_address,omitempty" bson:"email_address"`
	IsAdmin     string    `json:"is_admin,omitempty" bson:"is_admin"`
	Password    []byte    `json:"password,omitempty" bson:"password"`
}

// UserDTO transfers user data
type UserDTO struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	EmailAdress string `json:"email_address"`
	IsAdmin     string `json:"is_admin"`
	Password    string `json:"password"`
}
