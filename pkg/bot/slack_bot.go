package bot

import (
	"context"
	"magento/bot/pkg/repository"

	"github.com/slack-go/slack"
)

const (
	SlackToken    = "bot/slack/token"
	SlackChanelId = "bot/slack/chanel_id"
)

type SlackBot struct {
	app      *slack.Client
	chanelId string
}

func CreateSlackBot(cfgRepository repository.ConfigRepositoryInterface) (Bot, error) {
	token, err := cfgRepository.GetByPath(SlackToken, context.Background())
	if err != nil {
		return nil, err
	}
	chanelId, err := cfgRepository.GetByPath(SlackChanelId, context.Background())
	if err != nil {
		return nil, err
	}

	app := slack.New(token.Value)

	return SlackBot{app: app, chanelId: chanelId.Value}, nil
}

func (b SlackBot) SendMessage(message string) (bool, error) {
	_, _, _, err := b.app.SendMessage(b.chanelId, slack.MsgOptionText(message, false))
	if err != nil {
		return false, err
	}
	return true, err

}
