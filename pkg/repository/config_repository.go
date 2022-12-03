package repository

import (
	"fmt"
	"magento/bot/pkg/database"
	"magento/bot/pkg/model"

	"github.com/sirupsen/logrus"
)

type ConfigRepository struct {
	client database.PostgressConfigInterface
}

func NewConfigRepository(client database.PostgressConfigInterface) ConfigRepositoryInterface {
	return &ConfigRepository{client: client}
}

//get document by id
func (r *ConfigRepository) GetByPath(path string) (*model.Config, error) {
	var con model.Config
	row := r.client.GetByPath(path)
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(&con.Id, &con.Path, &con.Value)
	if con.Id == 0 {
		return &con, fmt.Errorf("there no config with path %s", path)
	}

	return &con, nil
}

//update config
//return true if ok and false and error
func (r *ConfigRepository) UpdateConfig(con *model.Config) (bool, error) {
	rowsAffected, err := r.client.UpdateById(con.Id, con.Path, con.Value)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("config with id %d doesn't update", con.Id)
	}
	return true, nil
}

//delete config
//return true if ok and false and error
func (r *ConfigRepository) DeleteConfig(id int64) (bool, error) {
	rowsAffected, err := r.client.DeleteById(id)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("config with id %d doesn't deleted", id)
	}
	return true, nil
}

//create config
//return true if ok and false and error
func (r *ConfigRepository) CreateConfig(con *model.Config) (int64, error) {
	id, err := r.client.Insert(con.Path, con.Value)
	if err != nil && id < 1 {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("config doesn't created")
	}
	return id, nil
}
