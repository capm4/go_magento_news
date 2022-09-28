FROM golang:latest

COPY ./ /bot
WORKDIR /bot
RUN go build -o bin/bot cmd/bot/bot.go
CMD ["bin/bot"]