package repository

import (
	"context"
	"fmt"
	"magento/bot/pkg/database"
	"magento/bot/pkg/model"

	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	client database.PostgressUserInterface
}

func NewUserRepository(client database.PostgressUserInterface) UserRepositoryInterface {
	return &UserRepository{client: client}
}

//get user by loger
func (r *UserRepository) GetByLogin(login string, ctx context.Context) (*model.User, error) {
	var user model.User
	row, err := r.client.GetByLogin(login, ctx)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(&user.Id, &user.Name, &user.Login, &user.Password, &user.UserRole, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if user.Login == "" {
		return &user, fmt.Errorf("there no user with login %s", login)
	}
	if !user.IsActive {
		return &user, fmt.Errorf("user with login %s is not active", login)
	}

	return &user, nil
}

//check if user exist by login
func (r *UserRepository) IsExist(login string, ctx context.Context) (bool, error) {
	var user model.User
	row, err := r.client.GetByLogin(login, ctx)
	if err != nil {
		return false, err
	}
	if row.Err() != nil {
		return false, row.Err()
	}
	row.Scan(&user.Id, &user.Name, &user.Login, &user.Password, &user.IsActive, &user.UserRole, &user.UpdatedAt, &user.CreatedAt)

	return user.Id != 0, nil
}

//update user
//return true if ok and false and error
func (r *UserRepository) UpdateUser(user *model.User, ctx context.Context) (bool, error) {
	rowsAffected, err := r.client.Update(*user, ctx)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("user with id %d doesn't update", user.Id)
	}
	return true, nil
}

//create user
//return true if ok and false and error
func (r *UserRepository) CreateUser(user *model.User, ctx context.Context) (int64, error) {
	id, err := r.client.Insert(*user, ctx)
	if err != nil && id < 1 {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("user doesn't created")
	}
	return id, nil
}
