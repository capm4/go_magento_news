DBconnection := $(shell cat config/DBconnectio.txt)

.PHONY: migrateup runDocker migratedown
.SILENT:

runDocker:
	docker-compose up -d

build: runDocker

migrateup: 
	migrate -path db/migration -database $(DBconnection) -verbose up

migratedown: 
	migrate -path db/migration -database "$(DBconnection)" -verbose down
