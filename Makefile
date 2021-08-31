build-pi:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o ./cmd/discordBot-Golang/ ./cmd/discordBot-Golang/*.go

build:
	go build -o ./cmd/discordBot-Golang/ ./cmd/discordBot-Golang/*.go

run:
	go run ./cmd/discordBot-Golang/

install: build
	ssh gamecenter "systemctl stop musicbot.service"
	scp cmd/discordBot-Golang/main gamecenter:/home/alcmoraes/Bots/discord/music/bot
	scp .env gamecenter:/home/alcmoraes/Bots/discord/music/.env
	ssh gamecenter "systemctl start musicbot.service"