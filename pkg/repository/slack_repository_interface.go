package repository

import (
	"context"
	"magento/bot/pkg/model"
)

type SlackRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*model.SlackBot, error)
	GetById(id int64, ctx context.Context) (*model.SlackBot, error)
	Update(slack *model.SlackBot, ctx context.Context) (bool, error)
	Create(slack *model.SlackBot, ctx context.Context) (int64, error)
	Delete(id int64, ctx context.Context) (bool, error)
	IsExistById(id int64, ctx context.Context) (bool, error)
	IsExistWebsiteInSlack(slackId, websiteId int64, ctx context.Context) (bool, error)
	InsertWebsiteToSlack(slackId, websiteId int64, ctx context.Context) (int64, error)
	RemoveWebsiteFromSlack(slackId, websiteId int64, ctx context.Context) (bool, error)
	GetAllWebsiteBySlackId(id int64, ctx context.Context) ([]*model.Website, error)
}
