package repository

import (
	"context"
	"magento/bot/pkg/model"
)

type ConfigRepositoryInterface interface {
	GetByPath(path string, ctx context.Context) (*model.Config, error)
	UpdateConfig(doc *model.Config, ctx context.Context) (bool, error)
	DeleteConfig(id int64, ctx context.Context) (bool, error)
	CreateConfig(doc *model.Config, ctx context.Context) (int64, error)
}
