package repository

import "magento/bot/pkg/model"

type ConfigRepositoryInterface interface {
	GetByPath(path string) (*model.Config, error)
	UpdateConfig(doc *model.Config) (bool, error)
	DeleteConfig(id int64) (bool, error)
	CreateConfig(doc *model.Config) (int64, error)
}
