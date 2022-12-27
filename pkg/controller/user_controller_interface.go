package controller

import (
	"github.com/labstack/echo"
)

type UserControllerInterface interface {
	AuthUser(c echo.Context, secret string) error
	Create(c echo.Context) error
	DeleteById(c echo.Context) error
	Update(c echo.Context) error
	GetUsers(c echo.Context) error
	GetUserByLogin(c echo.Context) error
}
