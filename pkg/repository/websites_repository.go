package repository

import (
	"context"
	"fmt"
	"magento/bot/pkg/database"
	"magento/bot/pkg/model"

	"github.com/sirupsen/logrus"
)

const domainNameWebsite = "website"

type WebsiteRepository struct {
	client database.PostgressWebsitesInterface
}

func NewWebsiteRepository(client database.PostgressWebsitesInterface) WebsiteRepositoryInterface {
	return &WebsiteRepository{client: client}
}

//get all websites from DB
func (r *WebsiteRepository) GetAll(ctx context.Context) ([]*model.Website, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	websites := []*model.Website{}
	rows, err := r.client.GetAll(c)
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
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	var doc model.Website
	row, err := r.client.GetById(id, c)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(&doc.Id, &doc.Url, &doc.Selector, &doc.Attribute, &doc.LastUrl)

	return &doc, nil
}

//update website
//return true if ok and false and error
func (r *WebsiteRepository) Update(website *model.Website, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rowsAffected, err := r.client.Update(*website, c)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("%s with id %d doesn't update", domainNameWebsite, website.Id)
	}
	return true, nil
}

//delete website
//return true if ok and false and error
func (r *WebsiteRepository) Delete(id int64, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rowsAffected, err := r.client.DeleteById(id, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("%s with id %d doesn't deleted", domainNameWebsite, id)
	}
	if rowsAffected < 1 {
		return false, fmt.Errorf("%s with id %d doesn't deleted", domainNameWebsite, id)
	}
	return true, nil
}

//create website
//return true if ok and false and error
func (r *WebsiteRepository) Create(website *model.Website, ctx context.Context) (int64, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	id, err := r.client.Insert(*website, c)
	if err != nil && id < 1 {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("%s doesn't created", domainNameWebsite)
	}
	return id, nil
}

//check by if if website exist
func (r *WebsiteRepository) IsExistById(id int64, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	exist, err := r.client.IsExistById(id, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("something goes wrong while checking if %s exist", domainNameWebsite)
	}
	return exist, nil
}
