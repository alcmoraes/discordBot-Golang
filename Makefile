build-pi:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o ./cmd/discordBot-Golang/ ./cmd/discordBot-Golang/*.go

build:
	go build -o ./cmd/discordBot-Golang/ ./cmd/discordBot-Golang/*.go

run:
	go run ./cmd/discordBot-Golang/

install: build
	scp cmd/discordBot-Golang/main alcmoraes@192.168.0.3:/home/alcmoraes/Bots/discord/music/bot
	scp .env alcmoraes@192.168.0.3:/home/alcmoraes/Bots/discord/music/.env