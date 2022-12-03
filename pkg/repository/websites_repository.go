package repository

import (
	"context"
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

//get all websites from DB
func (r *WebsiteRepository) GetAll(ctx context.Context) ([]*model.Website, error) {
	websites := []*model.Website{}
	rows, err := r.client.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		website := model.Website{}
		err = rows.Scan(&website.Id, &website.Url, &website.Selector, &website.Attribute, &website.LastUrl)
		if err != nil {
			logrus.Warning(err.Error())
		} else {
			websites = append(websites, &website)
		}
	}

	return websites, nil
}

//get document by id
func (r *WebsiteRepository) GetById(id int64, ctx context.Context) (*model.Website, error) {
	var doc model.Website
	row, err := r.client.GetById(id, ctx)
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
func (r *WebsiteRepository) Update(website *model.Website, ctx context.Context) (bool, error) {
	rowsAffected, err := r.client.Update(*website, ctx)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("website with id %d doesn't update", website.Id)
	}
	return true, nil
}

//delete website
//return true if ok and false and error
func (r *WebsiteRepository) Delete(id int64, ctx context.Context) (bool, error) {
	rowsAffected, err := r.client.DeleteById(id, ctx)
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
func (r *WebsiteRepository) Create(website *model.Website, ctx context.Context) (int64, error) {
	id, err := r.client.Insert(*website, ctx)
	if err != nil && id < 1 {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("website doesn't created")
	}
	return id, nil
}
