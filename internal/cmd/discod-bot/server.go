package cmd

import (
	"context"
	"log"

	commandsProvider "discordbot-golang/internal/commands/provider"
	"discordbot-golang/internal/discord"
	"discordbot-golang/internal/logger"
	messageProvider "discordbot-golang/internal/messages/provider"
	musicProvider "discordbot-golang/internal/music/provider"
	voiceProvider "discordbot-golang/internal/voice/provider"

	"discordbot-golang/internal/routes"
	youtubeProvider "discordbot-golang/internal/youtube/provider"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var (
	botToken string
)

func registerHooks(lifecycle fx.Lifecycle, discord discord.Discord) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				log.Print("Starting server.")
				if err := discord.OpenConnection(); err != nil {
					log.Printf("%v\n", err)
				}
				return nil
			},
			OnStop: func(context.Context) error {
				log.Print("Stopping server.")
				if err := discord.CloseConnection(); err != nil {
					log.Printf("%v\n", err)
				}
				return nil
			},
		},
	)
}

// RunServer runs discord bot server
func RunServer() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("dotEnv: can't loading .env file")
	}

	app := fx.New(
		fx.Provide(logger.NewLogger),
		fx.Provide(discord.NewSession),
		fx.Invoke(registerHooks),
		messageProvider.RepositoryModule,
		messageProvider.UsecaseModule,
		voiceProvider.UsecaseModule,
		youtubeProvider.UsecaseModule,
		musicProvider.UsecaseModule,
		messageProvider.DeliveryModule,
		commandsProvider.DeliveryModule,
		routes.Module,
	)
	app.Run()

	return nil
}
