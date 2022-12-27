package repository

import (
	"context"
	"magento/bot/pkg/model"
)

type UserRepositoryInterface interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByLogin(login string, ctx context.Context) (*model.User, error)
	IsExist(login string, ctx context.Context) (bool, error)
	UpdateUser(user *model.User, ctx context.Context) (bool, error)
	CreateUser(user *model.User, ctx context.Context) (int64, error)
	Delete(id int64, ctx context.Context) (bool, error)
	IsExistById(id int64, ctx context.Context) (bool, error)
}
