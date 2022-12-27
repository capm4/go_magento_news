package controller

import (
	"errors"
	"fmt"
	"magento/bot/pkg/model"
	"magento/bot/pkg/repository"
	"magento/bot/pkg/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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

func (uc *UserController) Create(c echo.Context) error {
	u, _ := GetUserFromContext(c)
	if u.Login == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("login is required"))
	}
	newUser, err := model.CreateUser(u.Name, u.Login, u.Password, model.RoleAdmin)
	if err != nil {
		logrus.Warning(err)
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while creating user"))
	}
	isExist, err := uc.userRepository.IsExist(newUser.Login, c.Request().Context())
	if err != nil {
		logrus.Warning(err)
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while creating user"))
	}
	if isExist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("user with login %s already exit", u.Login)))
	}
	id, err := uc.userRepository.CreateUser(newUser, c.Request().Context())
	if err != nil {
		logrus.Warning(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("user doesn't created"))
	}
	newUser.Id = id
	return c.JSON(http.StatusCreated, newUser.ToResponseUser())
}

func (us *UserController) DeleteById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while deleting user"))
	}
	ok, err := us.userRepository.Delete(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while delete user."))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg(fmt.Sprintf("user with id %d was deleted", id)))
}

func (d *UserController) Update(c echo.Context) error {
	w := &model.User{}
	err := w.Update(c)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating user"))
	}
	exist, err := d.userRepository.IsExistById(w.Id, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while checking if user exist"))
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("user with %d doesn't exist", w.Id)))
	}
	ok, err := d.userRepository.UpdateUser(w, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating user"))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while updating user"))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg("user was updated"))
}

func (u *UserController) GetUsers(c echo.Context) error {
	users, err := u.userRepository.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, model.ToResponseUsers(users))
}

func (u *UserController) GetUserByLogin(c echo.Context) error {
	login := c.Param("id")
	if login == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("login is required param"))
	}
	user, err := u.userRepository.GetByLogin(login, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, user.ToResponseUser())
}
