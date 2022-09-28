package config

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const EnvFileName string = ".env"
const DebugMode bool = true

var config *Сonfig

type Сonfig struct {
	SlackToken    string `env:"SlackToken"`
	SlackChanelId string `env:"SlackChanelId"`
	FilePath      string `env:"FilePath" envDefault:"src.json"`
	LogLevel      string `env:"LogLevel" envDefault:"info"`
	LogFilePath   string `env:"LogFilePath" envDefault:""`
	LogFormat     string `env:"LogFormat" envDefault:"text"`
	Debug         bool   `env:"Debug" envDefault:"false"`
}

func LoadConfigVariables(fileName string) (*Сonfig, error) {
	cnf := &Сonfig{}
	envMap, err := godotenv.Read(fileName)
	if err != nil {
		return cnf, errors.New("failed to load .env file")
	}
	if err := env.Parse(cnf, env.Options{Environment: envMap}); err != nil {
		return cnf, errors.New("failed to parse .env file to configuration variables")
	}

	return cnf, nil
}
