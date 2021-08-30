package routes

import (
	commandsDelivery "discordbot-golang/internal/commands/delivery"
	"discordbot-golang/internal/discord"
	messageDelivery "discordbot-golang/internal/messages/delivery"
	"go.uber.org/fx"
)

//NewRoutes new Routes Handler
func NewRoutes(discord discord.Discord, messageDelivery messageDelivery.Delivery, commandsDelivery commandsDelivery.Delivery) {
	discord.AddHandler(messageDelivery.GetMessageHandler)
	discord.AddHandler(commandsDelivery.GetCommandsHandler)
}

//Module .
var Module = fx.Options(
	fx.Invoke(NewRoutes),
)
