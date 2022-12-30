package router

import (
	"errors"
	"magento/bot/pkg/config"
	"magento/bot/pkg/controller"
	"magento/bot/pkg/model"
	"magento/bot/pkg/services"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	version1 = "/api/v1"
)

func NewRouter(e *echo.Echo, c controller.AppController, cfg *config.Сonfig) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	userRouter(e, c, cfg)
	websiteRouter(e, c, cfg)
	slackRouter(e, c, cfg)
	return e
}

func websiteRouter(e *echo.Echo, c controller.AppController, cfg *config.Сonfig) {
	v1JwtWithRole := e.Group(version1)
	addJWT(v1JwtWithRole, cfg)
	addJwtRoleCheck(v1JwtWithRole)
	v1JwtWithRole.GET("/websites", func(context echo.Context) error { return c.Website.GetWebsites(context) })
	v1JwtWithRole.GET("/website/:id", func(context echo.Context) error { return c.Website.GetWebsitesById(context) })
	v1JwtWithRole.PUT("/website", func(context echo.Context) error { return c.Website.Create(context) })
	v1JwtWithRole.DELETE("/website/:id", func(context echo.Context) error { return c.Website.DeleteById(context) })
	v1JwtWithRole.POST("/website", func(context echo.Context) error { return c.Website.Update(context) })
	v1JwtWithRole.POST("/check_website", func(context echo.Context) error { return c.Website.CheckWebsite(context) })
}

func userRouter(e *echo.Echo, c controller.AppController, cfg *config.Сonfig) {
	v1 := e.Group(version1)
	v1JwtWithRole := e.Group(version1)
	addJWT(v1JwtWithRole, cfg)
	addJwtRoleCheck(v1JwtWithRole)
	v1.POST("/user/login", func(context echo.Context) error { return c.User.AuthUser(context, c.Config.JwtSecret) })
	v1JwtWithRole.PUT("/user", func(context echo.Context) error { return c.User.Create(context) })
	v1JwtWithRole.POST("/user", func(context echo.Context) error { return c.User.Update(context) })
	v1JwtWithRole.DELETE("/user/:id", func(context echo.Context) error { return c.User.DeleteById(context) })
	v1JwtWithRole.GET("/users", func(context echo.Context) error { return c.User.GetUsers(context) })
	v1JwtWithRole.GET("/user/:id", func(context echo.Context) error { return c.User.GetUserByLogin(context) })
}

func slackRouter(e *echo.Echo, c controller.AppController, cfg *config.Сonfig) {
	v1JwtWithRole := e.Group(version1)
	addJWT(v1JwtWithRole, cfg)
	addJwtRoleCheck(v1JwtWithRole)
	v1JwtWithRole.GET("/slackbots", func(context echo.Context) error { return c.Slack.GetAll(context) })
	v1JwtWithRole.GET("/slackbot/:id", func(context echo.Context) error { return c.Slack.GetById(context) })
	v1JwtWithRole.PUT("/slackbot", func(context echo.Context) error { return c.Slack.Create(context) })
	v1JwtWithRole.POST("/slackbot", func(context echo.Context) error { return c.Slack.Update(context) })
	v1JwtWithRole.DELETE("/slackbot/:id", func(context echo.Context) error { return c.Slack.DeleteById(context) })
	v1JwtWithRole.POST("/slackbot/website", func(context echo.Context) error { return c.Slack.AddWebsiteToSlack(context) })
	v1JwtWithRole.DELETE("/slackbot/website", func(context echo.Context) error { return c.Slack.RemoveWebsitefromSlack(context) })
	v1JwtWithRole.GET("/slackbot/website/:id", func(context echo.Context) error { return c.Slack.GetAllWebsiteBySlackId(context) })
}

func addJWT(g *echo.Group, config *config.Сonfig) {
	jwtConfig := middleware.JWTConfig{
		SigningKey: []byte(config.JwtSecret),
		Claims:     &services.TokenClaims{},
	}
	g.Use(middleware.JWTWithConfig(jwtConfig))
}

func addJwtRoleCheck(g *echo.Group) {
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*services.TokenClaims)
			role := claims.UserRole
			if role != model.RoleAdmin {
				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "you don't have permission",
					Internal: errors.New("you don't have permission"),
				}
			}
			return next(c)
		}
	})
}
