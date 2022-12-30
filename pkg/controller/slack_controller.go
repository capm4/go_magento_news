package controller

import (
	"encoding/json"
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

func (s *SlackController) GetById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if idParama == "" {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while getting website"))
	}
	slack, err := s.slackRepository.GetById(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if slack.Id == 0 {
		return c.JSON(http.StatusNoContent, slack)
	}
	return c.JSON(http.StatusOK, slack)
}

func (s *SlackController) Create(c echo.Context) error {
	slack, err := model.CreateSlackFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	id, err := s.slackRepository.Create(slack, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	slack.Id = id
	return c.JSON(http.StatusCreated, slack)

}

func (s *SlackController) Update(c echo.Context) error {
	slack := &model.SlackBot{}
	err := s.Update(c)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating slack"))
	}
	exist, err := s.slackRepository.IsExistById(slack.Id, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while checking if slackbot exist"))
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("slackbot with %d doesn't exist", slack.Id)))
	}
	ok, err := s.slackRepository.Update(slack, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating slackbot"))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while updating slackbot"))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg("slackbot was updated"))
}

func (s *SlackController) DeleteById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while deleting slackbot"))
	}
	ok, err := s.slackRepository.Delete(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while delete slackbot."))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg(fmt.Sprintf("slackbot with id %d was deleted", id)))
}

func (s *SlackController) AddWebsiteToSlack(c echo.Context) error {
	data := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	keys := []string{"slackId", "websiteId"}
	parsetData, err := validateJsonAndParse(data, keys)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	exist, err := s.slackRepository.IsExistWebsiteInSlack(parsetData["slackId"], parsetData["websiteId"], c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if exist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("webiste with %d alredy in slack %d", parsetData["websiteId"], parsetData["slackId"])))
	}
	id, err := s.slackRepository.InsertWebsiteToSlack(parsetData["slackId"], parsetData["websiteId"], c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusCreated, CreateResponseMsg(fmt.Sprintf("website %d assign to slack %d, return id %d", parsetData["websiteId"], parsetData["slackId"], id)))
}

func (s *SlackController) RemoveWebsitefromSlack(c echo.Context) error {
	data := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	keys := []string{"slackId", "websiteId"}
	parsetData, err := validateJsonAndParse(data, keys)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	exist, err := s.slackRepository.IsExistWebsiteInSlack(parsetData["slackId"], parsetData["websiteId"], c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("webiste with %d not in slack %d", parsetData["websiteId"], parsetData["slackId"])))
	}
	ok, err := s.slackRepository.RemoveWebsiteFromSlack(parsetData["slackId"], parsetData["websiteId"], c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while delete website from slackbot."))
	}
	return c.JSON(http.StatusCreated, CreateResponseMsg(fmt.Sprintf("website %d removed from slack %d", parsetData["websiteId"], parsetData["slackId"])))
}

//validate json map and parse string to int
func validateJsonAndParse(data map[string]string, keys []string) (map[string]int64, error) {
	parsetData := make(map[string]int64)
	for _, key := range keys {
		id, ok := data[key]
		if !ok {
			return nil, fmt.Errorf("key %s doesn't exist", key)
		}
		value, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		parsetData[key] = value
	}
	return parsetData, nil
}

func (s *SlackController) GetAllWebsiteBySlackId(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while getting website by slackbot id"))
	}
	websites, err := s.slackRepository.GetAllWebsiteBySlackId(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseWebsites{Websites: websites})
}
