package database

import (
	"database/sql"
)

type PostgressConfigInterface interface {
	GetByPath(path string) *sql.Row
	UpdateById(id int64, path, value string) (int64, error)
	DeleteById(id int64) (int64, error)
	Insert(path, value string) (int64, error)
}
