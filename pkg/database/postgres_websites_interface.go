package database

import (
	"context"
	"database/sql"
	"magento/bot/pkg/model"
)

type PostgressWebsitesInterface interface {
	GetAll(ctx context.Context) (*sql.Rows, error)
	GetById(id int64, ctx context.Context) (*sql.Row, error)
	Update(website model.Website, ctx context.Context) (int64, error)
	DeleteById(id int64, ctx context.Context) (int64, error)
	Insert(website model.Website, ctx context.Context) (int64, error)
}
