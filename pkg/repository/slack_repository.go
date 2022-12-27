package repository

import (
	"context"
	"fmt"
	"magento/bot/pkg/database"
	"magento/bot/pkg/model"

	"github.com/sirupsen/logrus"
)

const domainNameSlackBot = "slack"

type SlackRepository struct {
	client database.PostgressSlackInterface
}

func NewSlackRepository(client database.PostgressSlackInterface) SlackRepositoryInterface {
	return &SlackRepository{client: client}
}

//get all slack bot from DB
func (r *SlackRepository) GetAll(ctx context.Context) ([]*model.SlackBot, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	slackBots := []*model.SlackBot{}
	rows, err := r.client.GetAll(c)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		slackBot := model.SlackBot{}
		err = rows.Scan(&slackBot.Id, &slackBot.Name, &slackBot.Token, &slackBot.ChannelId, &slackBot.CronEvery, &slackBot.LastCronRun)
		if err != nil {
			logrus.Warning(err.Error())
		} else {
			slackBots = append(slackBots, &slackBot)
		}
	}

	return slackBots, nil
}

//get slack bot by id
func (r *SlackRepository) GetById(id int64, ctx context.Context) (*model.SlackBot, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	var slackBot model.SlackBot
	row, err := r.client.GetById(id, c)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	row.Scan(&slackBot.Id, &slackBot.Name, &slackBot.Token, &slackBot.ChannelId, &slackBot.CronEvery, &slackBot.LastCronRun)

	return &slackBot, nil
}

//create slack bot
//return true if ok and false and error
func (r *SlackRepository) Create(slack *model.SlackBot, ctx context.Context) (int64, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	id, err := r.client.Insert(*slack, c)
	if err != nil && id < 1 {
		logrus.Warning(err.Error())
		return 0, fmt.Errorf("%s doesn't created", domainNameSlackBot)
	}
	return id, nil
}

//check by id if slackbot exist
func (r *SlackRepository) IsExistById(id int64, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	exist, err := r.client.IsExistById(id, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("something goes wrong while checking if %s exist", domainNameSlackBot)
	}
	return exist, nil
}

//update slackbot
//return true if ok and false and error
func (r *SlackRepository) Update(slackBot *model.SlackBot, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rowsAffected, err := r.client.Update(*slackBot, c)
	if err != nil && rowsAffected < 1 {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("%s with id %d doesn't update", domainNameSlackBot, slackBot.Id)
	}
	return true, nil
}

//delete slackbot
//return true if ok and false and error
func (r *SlackRepository) Delete(id int64, ctx context.Context) (bool, error) {
	c, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()
	rowsAffected, err := r.client.DeleteById(id, c)
	if err != nil {
		logrus.Warning(err.Error())
		return false, fmt.Errorf("%s with id %d doesn't deleted", domainNameSlackBot, id)
	}
	if rowsAffected < 1 {
		return false, fmt.Errorf("%s with id %d doesn't deleted", domainNameSlackBot, id)
	}
	return true, nil
}
