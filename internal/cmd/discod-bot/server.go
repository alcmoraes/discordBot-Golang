package cmd

import (
	"context"
	"log"

	commandsProvider "discordbot-golang/internal/commands/provider"
	"discordbot-golang/internal/discord"
	jackettProvider "discordbot-golang/internal/jackett/provider"
	"discordbot-golang/internal/logger"
	musicProvider "discordbot-golang/internal/music/provider"
	voiceProvider "discordbot-golang/internal/voice/provider"

	"discordbot-golang/internal/routes"
	youtubeProvider "discordbot-golang/internal/youtube/provider"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
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
		voiceProvider.UsecaseModule,
		jackettProvider.UsecaseModule,
		youtubeProvider.UsecaseModule,
		musicProvider.UsecaseModule,
		commandsProvider.DeliveryModule,
		routes.Module,
	)
	app.Run()

	return nil
}
