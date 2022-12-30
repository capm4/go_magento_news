package controller

import (
	"github.com/labstack/echo"
)

type WebsiteControllerInterface interface {
	GetWebsites(c echo.Context) error
	GetWebsitesById(c echo.Context) error
	Create(c echo.Context) error
	DeleteById(c echo.Context) error
	Update(c echo.Context) error
	CheckWebsite(c echo.Context) error
}
