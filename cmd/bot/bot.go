package main

import (
	"context"
	"fmt"
	"magento/bot/pkg/bot"
	"magento/bot/pkg/config"
	"magento/bot/pkg/logger"
	"magento/bot/pkg/registry"
	"magento/bot/pkg/router"
	"magento/bot/pkg/worker"
	"os"
	"os/signal"

	"github.com/labstack/echo"
	"gopkg.in/robfig/cron.v2"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfigVariables(config.EnvFileName)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	logger.InitLogger(cfg)
	registry, err := registry.Init(cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	c := cron.New()
	//c.AddFunc("0 */1 * * * *", func() { logrus.Info("ping") })
	//c.AddFunc("30 7 * * *", func() { RunBot(registry) })
	//c.Start()
	//RunBot(registry)
	//c.AddFunc("* * * * *", func() { RunBot(registry) })
	c.Start()
	RunRouter(registry)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func RunBot(registry *registry.Registry) {
	fmt.Println("ping")
	bot, err := bot.CreateSlackBot(registry.CfgRepository)
	//bot.SendMessage("ping")
	fmt.Println("ping")
	if err != nil {
		logrus.Warn(err)
	}
	websites, err := registry.WebRepository.GetAll(context.Background())
	if err != nil {
		logrus.Warn(err)
	}

	if bot != nil && websites != nil {
		for _, website := range websites {
			links := worker.GetLinks(*website)
			if len(links) > 0 {
				lastLink := links[0]
				website.LastUrl = lastLink
				registry.WebRepository.Update(website, context.Background())
				for _, link := range links {
					//bot.SendMessage(link)
					fmt.Println(link)
				}
			}
		}
	}
}

func RunRouter(registry *registry.Registry) {
	e := echo.New()
	controller, err := registry.NewAppController()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	e = router.NewRouter(e, controller, &registry.Config)
	fmt.Println("Server listen at http://localhost" + ":" + registry.Config.ServerAddressPort)
	if err := e.Start(":" + registry.Config.ServerAddressPort); err != nil {
		logrus.Fatalln(err)
	}
}
