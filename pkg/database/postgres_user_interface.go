package database

import (
	"context"
	"database/sql"
	"magento/bot/pkg/model"
)

type PostgressUserInterface interface {
	GetByLogin(login string, ctx context.Context) (*sql.Row, error)
	Update(user model.User, ctx context.Context) (int64, error)
	Insert(user model.User, ctx context.Context) (int64, error)
}
