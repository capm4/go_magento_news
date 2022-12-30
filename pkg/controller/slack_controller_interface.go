package controller

import (
	"github.com/labstack/echo"
)

type SlackControllerInterface interface {
	Create(c echo.Context) error
	DeleteById(c echo.Context) error
	Update(c echo.Context) error
	GetAll(c echo.Context) error
	GetById(c echo.Context) error
	AddWebsiteToSlack(c echo.Context) error
	RemoveWebsitefromSlack(c echo.Context) error
	GetAllWebsiteBySlackId(c echo.Context) error
}