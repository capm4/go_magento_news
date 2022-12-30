package database

import (
	"context"
	"database/sql"
	"magento/bot/pkg/model"
)

type PostgressSlackInterface interface {
	GetAll(ctx context.Context) (*sql.Rows, error)
	GetById(id int64, ctx context.Context) (*sql.Row, error)
	Update(website model.SlackBot, ctx context.Context) (int64, error)
	DeleteById(id int64, ctx context.Context) (int64, error)
	Insert(slack model.SlackBot, ctx context.Context) (int64, error)
	IsExistById(id int64, ctx context.Context) (bool, error)
	InsertWebsiteToSlack(slackId, websiteId int64, ctx context.Context) (int64, error)
	IsExistWebsiteInSlack(slackId, websiteId int64, ctx context.Context) (bool, error)
	DeleteWebsiteFromSlackById(slackId, websiteId int64, ctx context.Context) (int64, error)
	GetAllWebsiteBySlackId(id int64, ctx context.Context) (*sql.Rows, error)
}
