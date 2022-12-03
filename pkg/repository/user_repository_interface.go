package repository

import "magento/bot/pkg/model"

type UserRepositoryInterface interface {
	GetByLogin(login string) (*model.User, error)
	IsExist(login string) (bool, error)
	UpdateUser(user *model.User) (bool, error)
	CreateUser(user *model.User) (int64, error)
}
