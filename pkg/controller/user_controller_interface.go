package controller

import (
	"github.com/labstack/echo"
)

type UserControllerInterface interface {
	AuthUser(c echo.Context, secret string) error
}
