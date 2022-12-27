package worker

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"magento/bot/pkg/config"
	"os"
)

type Document struct {
	Url       string `json:"url"`
	Selector  string `json:"selector"`
	Attribute string `json:"attribute"`
	LastUrl   string `json:"last_url"`
}

func CreateDocuments(cfg *config.Сonfig) []Document {
	var documents []Document
	// Open our jsonFile
	jsonFile, err := os.Open(cfg.FilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		logrus.Warning(err)
		return documents
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	//unmarshal and make array of documents
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logrus.Warning(err)
		return documents
	}
	// we initialize our Users array
	err = json.Unmarshal(byteValue, &documents)
	if err != nil {
		logrus.Warning(err)
		return documents
	}
	return documents
}

func UpdateDocument(cfg *config.Сonfig, doc Document) {
	documents := CreateDocuments(cfg)
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
