package controller

import (
	"errors"
	"magento/bot/pkg/model"
	"magento/bot/pkg/repository"
	"magento/bot/pkg/services"
	"net/http"

	"github.com/labstack/echo"
)

type UserController struct {
	userRepository repository.UserRepositoryInterface
}

func NewUserController(userRepository repository.UserRepositoryInterface) (UserControllerInterface, error) {
	return &UserController{userRepository: userRepository}, nil
}

func (uc *UserController) AuthUser(c echo.Context, secret string) error {
	u, _ := GetUserFromContext(c)
	if u.Login == "" || u.Password == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("login and password are required param"))
	}
	c.Request().Context()
	user, err := uc.userRepository.GetByLogin(u.Login, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !services.CheckPasswordHash(u.Password, user.Password) {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("password is wrong"))
	}
	token, err := services.CreateJWT(secret, user.UserRole, user.Login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while generated token"))
	}
	return c.JSON(http.StatusOK, model.AuthResponse{Token: "Bearer " + token})
}

func GetUserFromContext(c echo.Context) (*model.User, error) {
	user := model.User{}
	if err := c.Bind(&user); !errors.Is(err, nil) {
		return nil, err
	}
	return &user, nil
}
