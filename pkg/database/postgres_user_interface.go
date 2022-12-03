package database

import (
	"database/sql"
	"magento/bot/pkg/model"
)

type PostgressUserInterface interface {
	GetByLogin(login string) (*sql.Row, error)
	Update(model.User) (int64, error)
	Insert(model.User) (int64, error)
}
