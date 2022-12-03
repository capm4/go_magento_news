package database

import (
	"context"
	"database/sql"
)

type PostgressConfigInterface interface {
	GetByPath(path string, ctx context.Context) *sql.Row
	UpdateById(id int64, path, value string, ctx context.Context) (int64, error)
	DeleteById(id int64, ctx context.Context) (int64, error)
	Insert(path, value string, ctx context.Context) (int64, error)
}
