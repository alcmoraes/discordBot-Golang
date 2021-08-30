package routes

import (
	commandsDelivery "discordbot-golang/internal/commands/delivery"
	"discordbot-golang/internal/discord"

	"go.uber.org/fx"
)

//NewRoutes new Routes Handler
func NewRoutes(discord discord.Discord, commandsDelivery commandsDelivery.Delivery) {
	discord.AddHandler(commandsDelivery.GetCommandsHandler)
}

//Module .
var Module = fx.Options(
	fx.Invoke(NewRoutes),
)
