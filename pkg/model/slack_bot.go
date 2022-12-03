package model

import (
	"time"

	"github.com/slack-go/slack"
)

type SlackBot struct {
	Id          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Token       string     `json:"token" db:"token"`
	ChannelId   string     `json:"channelId" db:"channel_id"`
	CronEvery   int64      `json:"CronEvery" db:"cron_every"`
	LastCronRun *time.Time `json:"LastCronRun" db:"last_cron_run"`
}

func (s *SlackBot) Send(msg string) (bool, error) {
	app := slack.New(s.Token)
	_, _, _, err := app.SendMessage(s.ChannelId, slack.MsgOptionText(msg, false))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SlackBot) IsRun() bool {
	currentTime := time.Now()
	cronTime := s.LastCronRun.Add(time.Minute * time.Duration(s.CronEvery))
	return currentTime.Sub(cronTime) >= 0
}
