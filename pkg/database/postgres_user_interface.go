package database

import (
	"context"
	"database/sql"
	"magento/bot/pkg/model"
)

type PostgressUserInterface interface {
	GetAll(ctx context.Context) (*sql.Rows, error)
	GetByLogin(login string, ctx context.Context) (*sql.Row, error)
	Update(user model.User, ctx context.Context) (int64, error)
	Insert(user model.User, ctx context.Context) (int64, error)
	IsExistByLogin(login string, ctx context.Context) (bool, error)
	DeleteById(id int64, ctx context.Context) (int64, error)
	IsExistById(id int64, ctx context.Context) (bool, error)
}
