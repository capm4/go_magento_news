package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo"
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

type ResponsSlackBot struct {
	SlackBot []*SlackBot `json:"slack_bot"`
}

func CreateSlackFromContext(ctx echo.Context) (*SlackBot, error) {
	slack := SlackBot{}
	err := slack.bindAndValidate(ctx)
	if err != nil {
		return nil, err
	}
	ct := time.Now()
	slack.LastCronRun = &ct
	return &slack, nil
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

//update and validate slack bot from echo context
func (s *SlackBot) Update(c echo.Context) error {
	err := s.bindAndValidate(c)
	if err != nil {
		return err
	}
	return nil
}

func (s *SlackBot) bindAndValidate(ctx echo.Context) error {
	if err := ctx.Bind(s); !errors.Is(err, nil) {
		return err
	}
	err := s.validate()
	if err != nil {
		return err
	}
	return nil
}

func (s *SlackBot) validate() error {
	if s.ChannelId == "" {
		return fmt.Errorf("chanel id is required")
	}
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}
