package provider

import (
	"discordbot-golang/internal/music/usecase"
	"go.uber.org/fx"
)

//UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewMusicUsecase),
)
