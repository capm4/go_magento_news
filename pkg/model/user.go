package model

import (
	"errors"
	"fmt"
	"magento/bot/pkg/services"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	RoleAdmin = "admin"
)

type User struct {
	Id        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Login     string     `json:"login" db:"login"`
	Password  string     `json:"password" db:"password"`
	IsActive  bool       `json:"is_active" db:"is_active" default:"true"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	UserRole  string     `json:"role" db:"role"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func CreateUser(name, login, password, user_role string) (*User, error) {
	user := User{}
	user.Name = name
	user.Login = login
	user.Password = password
	err := user.validate()
	if err != nil {
		return nil, err
	}
	if user_role != RoleAdmin {
		return nil, fmt.Errorf("incorect user role %s", user_role)
	}
	user.IsActive = true
	err = user.validateAndHasPassword()
	if err != nil {
		logrus.Warning(err.Error())
		return nil, err
	}
	if user.UserRole == "" {
		user.UserRole = RoleAdmin
	}
	return &user, nil
}

// future function for create user from REST API echo
// func (u *User) bindAndValidate(ctx echo.Context) error {
// 	if err := ctx.Bind(u); !errors.Is(err, nil) {
// 		return err
// 	}
// 	err := u.validate()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (u *User) validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Login == "" {
		return fmt.Errorf("login is required")
	}
	return nil
}

func (u *User) validateAndHasPassword() error {
	err := services.ValidatePassword(u.Password)
	if err != nil {
		return err
	}
	password, err := services.HashPassword(u.Password)
	if err != nil {
		logrus.Warning(err.Error())
		return errors.New("something when wrong while create new user")
	}
	u.Password = password
	return nil
}
