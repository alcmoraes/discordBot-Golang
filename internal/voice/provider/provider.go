package provider

import (
	"discordbot-golang/internal/voice/usecase"
	"go.uber.org/fx"
)

//UsecaseModule .
var UsecaseModule = fx.Options(
	fx.Provide(usecase.NewVoiceUsecase),
)
