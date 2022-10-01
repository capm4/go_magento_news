package main

import (
	"fmt"
	"magento/bot/pkg/bot"
	"magento/bot/pkg/config"
	"magento/bot/pkg/logger"
	"magento/bot/pkg/worker"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	cfg, err := config.LoadConfigVariables(config.EnvFileName)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logger.InitLogger(cfg)
	slack := bot.CreateSender(cfg)
	c := cron.New()

	c.AddFunc("0 */1 * * * *", func() { logrus.Info("ping") })
	c.AddFunc("30 7 * * *", func() { RunBot(cfg, slack) })
	c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func RunBot(cfg *config.Сonfig, slack *slack.Client) {
	docs := worker.CreateDocuments(cfg)
	for _, doc := range docs {
		links := worker.GetLinks(doc)
		if len(links) > 0 {
			lastLink := links[0]
			doc.LastUrl = lastLink
			worker.UpdateDocument(cfg, doc)
			for _, link := range links {
				bot.SendMessage(slack, cfg.SlackChanelId, link)
				fmt.Println(link)
			}
		}
	}
}
