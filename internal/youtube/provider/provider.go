package provider

import (
	"discordbot-golang/internal/youtube/usecase"
	"go.uber.org/fx"
)

//UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewYoutubeUsecase),
)
