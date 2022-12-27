package repository

import (
	"context"
	"fmt"
	"magento/bot/pkg/database"
	"magento/bot/pkg/model"

	"github.com/sirupsen/logrus"
)

const domainNameUser = "user"

type UserRepository struct {
	client database.PostgressUserInterface
}

func NewUserRepository(client database.PostgressUserInterface) UserRepositoryInterface {
	return &UserRepository{client: client}
}

//get all user from DB
func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	users := []*model.User{}
	rows, err := r.client.GetAll(c)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := model.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Login, &user.Password, &user.UserRole, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			logrus.Warning(err.Error())
		} else {
			users = append(users, &user)
		}
	}

	return users, nil
}

//get user by loger
func (r *UserRepository) GetByLogin(login string, ctx context.Context) (*model.User, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	var user model.User
	row, err := r.client.GetByLogin(login, c)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(&user.Id, &user.Name, &user.Login, &user.Password, &user.UserRole, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if user.Login == "" {
		return &user, fmt.Errorf("there no %s with login %s", domainNameUser, login)
	}
	if !user.IsActive {
		return &user, fmt.Errorf("%s with login %s is not active", domainNameUser, login)
	}

	return &user, nil
}

//check if user exist by login
func (r *UserRepository) IsExist(login string, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	exist, err := r.client.IsExistByLogin(login, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("something goes wrong while checking if %s exist", domainNameUser)
	}
	return exist, nil
}

//update user
//return true if ok and false and error
func (r *UserRepository) UpdateUser(user *model.User, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rowsAffected, err := r.client.Update(*user, c)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("%s with id %d doesn't update", domainNameUser, user.Id)
	}
	return true, nil
}

//create user
//return true if ok and false and error
func (r *UserRepository) CreateUser(user *model.User, ctx context.Context) (int64, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	id, err := r.client.Insert(*user, c)
	if err != nil {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("%s doesn't created", domainNameUser)
	}
	return id, nil
}

func (r *UserRepository) Delete(id int64, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rowsAffected, err := r.client.DeleteById(id, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("%s with id %d doesn't deleted", domainNameUser, id)
	}
	if rowsAffected < 1 {
		return false, fmt.Errorf("%s with id %d doesn't deleted", domainNameUser, id)
	}
	return true, nil
}

//check by if if user exist
func (r *UserRepository) IsExistById(id int64, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	exist, err := r.client.IsExistById(id, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("something goes wrong while checking if %s exist", domainNameUser)
	}
	return exist, nil
}
