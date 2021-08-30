package usecase

import (
	"fmt"
	"log"

	"discordbot-golang/internal/discord"
	voiceUsecase "discordbot-golang/internal/voice/usecase"
	youtubeUsecase "discordbot-golang/internal/youtube/usecase"

	"github.com/bwmarrin/discordgo"
)

//Usecase interface
type Usecase interface {
	PlayYoutubeURL(string, *discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
}

type musicUsecase struct {
	youtubeUsecase youtubeUsecase.Usecase
	voiceUsecase   voiceUsecase.Usecase
	discord        discord.Discord
}

//NewMusicUsecase new message delivery
func NewMusicUsecase(discord discord.Discord, yu youtubeUsecase.Usecase, vu voiceUsecase.Usecase) Usecase {
	return &musicUsecase{
		youtubeUsecase: yu,
		voiceUsecase:   vu,
		discord:        discord,
	}
}

func (mu musicUsecase) PlayYoutubeURL(url string, s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	voiceConnection, err := mu.voiceUsecase.ConnectToVoiceChannel(s, m, guild, true)
	if err != nil {
		log.Printf("Error: connect to voice channel, Message: '%s'", err)
		mu.discord.SendMessageToChannel(m.ChannelID, "Can't send message to channel. Something happened.")
		return
	}

	if discord.GetVoiceStatus() {
		mu.discord.SendMessageToChannel(m.ChannelID, "Wait for the music ends")
		return
	}
	youtubeInfo, err := mu.youtubeUsecase.GetYoutubeDownloadURL(url)
	if err != nil {
		log.Printf("Error: can't get youtube download url, Message: '%s'", err)
		mu.discord.SendMessageToChannel(m.ChannelID, "Can't find this music.")
		return
	}
	msg := fmt.Sprintf("Playing '%s' next", youtubeInfo.Title)
	mu.discord.SendMessageToChannel(m.ChannelID, msg)
	mu.voiceUsecase.PlayAudioFile(youtubeInfo.DownloadLink, voiceConnection)
}
