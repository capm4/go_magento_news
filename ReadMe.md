This application run check and sent news to slack chanel
=====================================
Logger and databases configuration: 
1) Add token to file .env

.ENV file configurations
=====================================
1) LogLevel = log write level. Options : trace, debug, info, warm, error, fatal, panic.
2) LogFilePath = path to logger file. Default var/log/system.log
3) LogFormat = log format. Has two options: text or json. Default text.
4) Debug = default true
5) DB_HOST = Postgres Working host
6) DB_PORT = Postgres Port
7) DB_NAME = Postgres database name
8) DB_USER = Postgres user
9) APP_PORT = App Port

migrate DB 
====================================
https://github.com/golang-migrate/migrate

generate customer by CLI
=====================================
1) go run cmd/cli/createUser.go 
 -l <login required> 
 -p <password required>
 -n <name required>  
 -r <role required (options :admin)>
 example go run cmd/cli/createUser.go -l=admin -p=admin123 -n=admin -r=admin

Project package:
=====================================
1) Entry point is in cmd/rest/main.go
2) magento/bot/pkg/config : response for application configuration
3) magento/bot/pkg/logger : module for configuration logger
4) magento/bot/pkg/model : package for domains
5) magento/bot/pkg/database : package for database 
6) magento/bot/pkg/bot : package for slack bot 
7) magento/bot/pkg/repository : repository to work with Postgres
8) magento/bot/pkg/registry : package for registry
9) magento/bot/pkg/service : package for services

Project Dependency
=====================================
1) Read env files : github.com/joho/godotenv
2) Logger https://github.com/sirupsen/logrus
3) Postgres database/sql
4) Slack bot github.com/slack-go/slack

Docker
======================================
1) postgres
2) golang:latest
