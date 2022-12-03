package repository

import (
	"fmt"
	"magento/bot/pkg/database"
	"magento/bot/pkg/model"

	"github.com/sirupsen/logrus"
)

type WebsiteRepository struct {
	client database.PostgressWebsitesInterface
}

func NewWebsiteRepository(client database.PostgressWebsitesInterface) WebsiteRepositoryInterface {
	return &WebsiteRepository{client: client}
}

//get all documents from DB
func (r *WebsiteRepository) GetAll() ([]*model.Website, error) {
	var documents []*model.Website
	rows, err := r.client.GetAll()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		document := model.Website{}
		err = rows.Scan(&document.Id, &document.Url, &document.Selector, &document.Attribute, &document.LastUrl)
		if err != nil {
			logrus.Warning(err.Error())
		} else {
			documents = append(documents, &document)
		}
	}

	return documents, nil
}

//get document by id
func (r *WebsiteRepository) GetById(id int64) (*model.Website, error) {
	var doc model.Website
	row, err := r.client.GetById(id)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(&doc.Id, &doc.Url, &doc.Selector, &doc.Attribute, &doc.LastUrl)
	if doc.Id == 0 {
		return &doc, fmt.Errorf("there no website with id %d", id)
	}

	return &doc, nil
}

//update website
//return true if ok and false and error
func (r *WebsiteRepository) Update(website *model.Website) (bool, error) {
	rowsAffected, err := r.client.Update(*website)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("website with id %d doesn't update", website.Id)
	}
	return true, nil
}

//delete website
//return true if ok and false and error
func (r *WebsiteRepository) Delete(id int64) (bool, error) {
	rowsAffected, err := r.client.DeleteById(id)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("website with id %d doesn't deleted", id)
	}
	if rowsAffected < 1 {
		return false, fmt.Errorf("website with id %d doesn't deleted", id)
	}
	return true, nil
}

//create website
//return true if ok and false and error
func (r *WebsiteRepository) Create(website *model.Website) (int64, error) {
	id, err := r.client.Insert(*website)
	if err != nil && id < 1 {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("website doesn't created")
	}
	return id, nil
}
