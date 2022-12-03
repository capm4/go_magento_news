package repository

import "magento/bot/pkg/model"

type WebsiteRepositoryInterface interface {
	GetAll() ([]*model.Website, error)
	GetById(id int64) (*model.Website, error)
	Update(doc *model.Website) (bool, error)
	Delete(id int64) (bool, error)
	Create(doc *model.Website) (int64, error)
}
