package database

import (
	"database/sql"
	"magento/bot/pkg/model"
)

type PostgressWebsitesInterface interface {
	GetAll() (*sql.Rows, error)
	GetById(id int64) (*sql.Row, error)
	Update(website model.Website) (int64, error)
	DeleteById(id int64) (int64, error)
	Insert(website model.Website) (int64, error)
}
