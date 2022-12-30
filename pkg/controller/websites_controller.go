package controller

import (
	"fmt"
	"magento/bot/pkg/model"
	"magento/bot/pkg/repository"
	"magento/bot/pkg/worker"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type WebsiteController struct {
	websiteRepository repository.WebsiteRepositoryInterface
}

func NewWebsiteController(docRepository repository.WebsiteRepositoryInterface) (WebsiteControllerInterface, error) {
	return &WebsiteController{websiteRepository: docRepository}, nil
}

func (d *WebsiteController) GetWebsites(c echo.Context) error {
	websites, err := d.websiteRepository.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, model.ResponseWebsites{Websites: websites})
}

func (d *WebsiteController) GetWebsitesById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if idParama == "" {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while getting website"))
	}
	website, err := d.websiteRepository.GetById(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, website.ToResponseWebsite())
}

func (d *WebsiteController) Create(c echo.Context) error {
	website, err := model.CreateWebsiteFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	id, err := d.websiteRepository.Create(website, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	website.Id = id
	return c.JSON(http.StatusCreated, website.ToResponseWebsite())

}

func (d *WebsiteController) DeleteById(c echo.Context) error {
	idParama := c.Param("id")
	if idParama == "" {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("id is required param"))
	}
	id, err := strconv.ParseInt(idParama, 10, 64)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while deleting website"))
	}
	ok, err := d.websiteRepository.Delete(id, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong while delete website."))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg(fmt.Sprintf("website with id %d was deleted", id)))
}

func (d *WebsiteController) Update(c echo.Context) error {
	w := &model.Website{}
	err := w.Update(c)
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating website"))
	}
	exist, err := d.websiteRepository.IsExistById(w.Id, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while checking if website exist"))
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(fmt.Sprintf("website with %d doesn't exist", w.Id)))
	}
	ok, err := d.websiteRepository.Update(w, c.Request().Context())
	if err != nil {
		logrus.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("error while updating website"))
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse("something goes wrong, while updating website"))
	}
	return c.JSON(http.StatusOK, CreateResponseMsg("website was updated"))
}

func (d *WebsiteController) CheckWebsite(c echo.Context) error {
	website, err := model.CreateWebsiteFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorResponse(err.Error()))
	}
	body := worker.GetLinks(*website)

	return c.JSON(http.StatusOK, body)
}
