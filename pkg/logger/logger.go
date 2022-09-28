package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"magento/bot/pkg/config"
	"os"
	"strings"
)

func InitLogger(cfg *config.Сonfig) error {
	err := setFormat(cfg.LogFormat)
	if err != nil {
		return err
	}
	err = setLogLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	if cfg.LogFilePath == "" {
		logrus.SetOutput(os.Stdout)
		return err
	}
	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("error opening file: %s", cfg.LogFilePath))
	}
	logrus.SetOutput(logFile)
	return nil
}

func setLogLevel(level string) error {
	levelType, err := logrus.ParseLevel(level)
	if err != nil {
		return errors.New(fmt.Sprintf("in env file set wrong log level. It should be one of (trace, debug, info, warm, error, fatal, panic) you set %s", level))
	}
	logrus.SetLevel(levelType)
	return nil
}

func setFormat(format string) error {
	switch format := strings.ToLower(format); format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		return errors.New(fmt.Sprintf("in env file set wrong log format. It should be json or text you set %s", format))
	}

	return nil
}
