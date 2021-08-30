package provider

import (
	"discordbot-golang/internal/messages/delivery"
	"discordbot-golang/internal/messages/repository"
	"discordbot-golang/internal/messages/usecase"
	"go.uber.org/fx"
)

//DeliveryModule .
var DeliveryModule = fx.Options(
	fx.Provide(delivery.NewMessageDelivery),
)

//RepositoryModule .
var RepositoryModule = fx.Options(
	fx.Provide(repository.NewMessageRepository),
)

//UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewMessagesUsecase),
)
