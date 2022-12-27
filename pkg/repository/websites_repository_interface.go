package repository

import (
	"context"
	"magento/bot/pkg/model"
)

type WebsiteRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*model.Website, error)
	GetById(id int64, ctx context.Context) (*model.Website, error)
	Update(doc *model.Website, ctx context.Context) (bool, error)
	Delete(id int64, ctx context.Context) (bool, error)
	Create(doc *model.Website, ctx context.Context) (int64, error)
	IsExistById(id int64, ctx context.Context) (bool, error)
}
