package provider

import (
	"discordbot-golang/internal/commands/delivery"
	"go.uber.org/fx"
)

//DeliveryModule .
var DeliveryModule = fx.Options(
	fx.Provide(delivery.NewCommandsDelivery),
)
