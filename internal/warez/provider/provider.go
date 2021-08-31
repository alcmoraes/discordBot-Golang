package provider

import (
	"discordbot-golang/internal/warez/usecase"

	"go.uber.org/fx"
)

//UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewWarezUsecase),
)
