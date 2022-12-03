package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"magento/bot/pkg/config"
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type Website struct {
	Id        int64  `json:"id" db:"id"`
	Url       string `json:"url" db:"url"`
	Selector  string `json:"selector" db:"selector"`
	Attribute string `json:"attribute" db:"attribute"`
	LastUrl   string `json:"last_url" db:"last_url"`
}

type ResponseWebsites struct {
	Websites []*Website `json:"websites"`
}

type ResponseWebsite struct {
	Id        int64  `json:"id" db:"id"`
	Url       string `json:"url" db:"url"`
	Selector  string `json:"selector" db:"selector"`
	Attribute string `json:"attribute" db:"attribute"`
	LastUrl   string `json:"last_url" db:"last_url"`
}

func CreateWebsite(cfg *config.Сonfig) []Website {
	var websites []Website
	// Open our jsonFile
	jsonFile, err := os.Open(cfg.FilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		logrus.Warning(err)
		return websites
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	//unmarshal and make array of documents
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logrus.Warning(err)
		return websites
	}
	// we initialize our Users array
	err = json.Unmarshal(byteValue, &websites)
	if err != nil {
		logrus.Warning(err)
		return websites
	}
	return websites
}

func CreateWebsiteFromContext(ctx echo.Context) (*Website, error) {
	website := Website{}
	err := website.bindAndValidate(ctx)
	if err != nil {
		return nil, err
	}
	return &website, nil
}

func UpdateDocument(cfg *config.Сonfig, doc Website) {
	documents := CreateWebsite(cfg)
	for i, document := range documents {
		if document.Url == doc.Url {
			documents[i].LastUrl = doc.LastUrl
			break
		}
	}
	jsonData, err := json.Marshal(documents)
	if err != nil {
		logrus.Warning(err)
	}
	err = ioutil.WriteFile(cfg.FilePath, jsonData, 0644)
	if err != nil {
		logrus.Warning(err)
	}
}

func (w *Website) ToResponseWebsite() *ResponseWebsite {
	return &ResponseWebsite{Id: w.Id, Url: w.Url, Selector: w.Selector, Attribute: w.Attribute, LastUrl: w.LastUrl}
}

func (w *Website) bindAndValidate(ctx echo.Context) error {
	if err := ctx.Bind(w); !errors.Is(err, nil) {
		return err
	}
	err := w.validate()
	if err != nil {
		return err
	}
	return nil
}

func (w *Website) validate() error {
	if w.Url == "" {
		return fmt.Errorf("url is required")
	}
	if w.Attribute == "" {
		return fmt.Errorf("attribute is required")
	}
	if w.Selector == "" {
		return fmt.Errorf("selector is required")
	}
	if w.LastUrl == "" {
		return fmt.Errorf("lastUrl is required")
	}
	return nil
}
