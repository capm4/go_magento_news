package bot

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"magento/bot/pkg/config"
)

func CreateSender(cfg *config.Сonfig) *slack.Client {
	return slack.New(cfg.SlackToken)
}

func SendMessage(app *slack.Client, chanelId string, message string) {
	_, _, _, err := app.SendMessage(chanelId, slack.MsgOptionText(message, false))
	if err != nil {
		logrus.Warning(err)
	}
}
