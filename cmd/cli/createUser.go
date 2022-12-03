package main

import (
	"flag"
	"fmt"
	"magento/bot/pkg/config"
	"magento/bot/pkg/logger"
	"magento/bot/pkg/model"
	"magento/bot/pkg/registry"
	"magento/bot/pkg/repository"
	"magento/bot/pkg/services"

	"github.com/sirupsen/logrus"
)

var userRepository repository.UserRepositoryInterface

func main() {
	login := flag.String("l", "", "login of user, required")
	password := flag.String("p", "", "password of user, required")
	name := flag.String("n", "", "name of user, required")
	role := flag.String("r", "", "role of user (admin), required")
	flag.Parse()
	if *login == "" || *password == "" || *name == "" || *role == "" {
		flag.Usage()
		return
	}
	if *role != model.RoleAdmin {
		flag.Usage()
		return
	}
	if checkIfUserExist(*login) {
		fmt.Printf("User with login %s already exist\n", *login)
		return
	}
	pas := createPasswordHash(*password)
	user := &model.User{
		Login:    *login,
		Name:     *name,
		Password: pas,
		IsActive: true,
		UserRole: *role,
	}
	saveUser(user)
	fmt.Printf("User with login %s successfuly created\n", *login)
}

func checkIfUserExist(login string) bool {
	repo := getUserRepository()
	isExist, err := repo.IsExist(login)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return isExist
}

func createPasswordHash(password string) string {
	err := services.ValidatePassword(password)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	pas, err := services.HashPassword(password)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return pas
}

func saveUser(user *model.User) {
	repo := getUserRepository()
	res, err := repo.CreateUser(user)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	if res < 0 {
		logrus.Fatal("user with login %s doesn't created", user.Login)
	}
}

func getUserRepository() repository.UserRepositoryInterface {
	if userRepository != nil {
		return userRepository
	}
	cfg, err := config.LoadConfigVariables(config.EnvFileName)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	err = logger.InitLogger(cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	db, err := registry.CreateDBConnection(cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	userRepository, err = registry.CreateUserRepository(*db)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return userRepository
}
