package models

import "golang.org/x/crypto/bcrypt"

// User represents the user model for the application.
// @Description User model
// @Name User
// @Property id int64 `json:"id,omitempty" bson:"id,omitempty"`
// @Property name string `json:"name,omitempty"`
// @Property email string `json:"email,omitempty"`
// @Property password string `json:"password,omitempty"`
type User struct {
	ID       int64  `json:"id,omitempty" bson:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateUserResponse struct {
    Message string `json:"message"`
    User    User   `json:"user"`
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
