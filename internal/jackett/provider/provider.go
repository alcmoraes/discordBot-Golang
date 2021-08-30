package provider

import (
	"discordbot-golang/internal/jackett/usecase"

	"go.uber.org/fx"
)

//UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewJackettUsecase),
)
