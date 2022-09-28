.PHONY:
.SILENT:

runDocker:
	docker-compose up -d

build: runDocker