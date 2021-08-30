package main

import (
	"log"
	"os"

	discordBot "discordbot-golang/internal/cmd/discod-bot"
)

func main() {
	if err := discordBot.RunServer(); err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	return
}
