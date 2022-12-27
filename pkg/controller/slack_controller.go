package controller

import (
	"fmt"
	"magento/bot/pkg/model"
	"magento/bot/pkg/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type SlackController struct {
	slackRepository repository.SlackRepositoryInterface
}

func NewSlackController(slackrepositore repository.SlackRepositoryInterface) (SlackControllerInterface, error) {
	return &SlackController{slackRepository: slackrepositore}, nil
}

func (s *SlackController) GetAll(c echo.Context) error {
	slackBots, err := s.slackRepository.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponsSlackBot{SlackBot: slackBots})
}

func (d *SlackController) GetById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if idParama == "" {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while getting website"))
	}
	slack, err := d.slackRepository.GetById(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if slack.Id == 0 {
		return c.JSON(http.StatusNoContent, slack)
	}
	return c.JSON(http.StatusOK, slack)
}

func (d *SlackController) Create(c echo.Context) error {
	slack, err := model.CreateSlackFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	id, err := d.slackRepository.Create(slack, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	slack.Id = id
	return c.JSON(http.StatusCreated, slack)

}

func (d *SlackController) Update(c echo.Context) error {
	s := &model.SlackBot{}
	err := s.Update(c)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating slack"))
	}
	exist, err := d.slackRepository.IsExistById(s.Id, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while checking if slackbot exist"))
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("slackbot with %d doesn't exist", s.Id)))
	}
	ok, err := d.slackRepository.Update(s, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating slackbot"))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while updating slackbot"))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg("slackbot was updated"))
}

func (d *SlackController) DeleteById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while deleting slackbot"))
	}
	ok, err := d.slackRepository.Delete(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while delete slackbot."))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg(fmt.Sprintf("slackbot with id %d was deleted", id)))
}
