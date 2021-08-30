package delivery

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordbot-golang/internal/discord"
	musicUsecase "discordbot-golang/internal/music/usecase"
	voiceUsecase "discordbot-golang/internal/voice/usecase"
)

//Delivery interface
type Delivery interface {
	GetCommandsHandler(*discordgo.Session, *discordgo.MessageCreate)
}

type commandsDelivery struct {
	musicUsecase musicUsecase.Usecase
	voiceUsecase voiceUsecase.Usecase
	discord      discord.Discord
}

//NewCommandsDelivery new message delivery
func NewCommandsDelivery(discord discord.Discord, mu musicUsecase.Usecase, vu voiceUsecase.Usecase) Delivery {
	return &commandsDelivery{
		musicUsecase: mu,
		voiceUsecase: vu,
		discord:      discord,
	}
}

func (cd commandsDelivery) GetCommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	botPrefix := os.Getenv("BOT_PREFIX")
	if botPrefix == "" {
		botPrefix = "~"
	}
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		log.Println(err)
	}

	if !strings.HasPrefix(m.Content, botPrefix) {
		return
	}

	if strings.HasPrefix(m.Content, botPrefix+"help") {
		help := fmt.Sprintf("**Help**\n==============================\n`%splay [Youtube Link]` : Plays some music\n`%sstop` : Stop playing\n`%sjoin` : Joins your channel\n==============================", botPrefix, botPrefix, botPrefix)
		cd.discord.SendMessageToChannel(m.ChannelID, help)
	} else if strings.HasPrefix(m.Content, botPrefix+"join") {
		cd.voiceUsecase.ConnectToVoiceChannel(s, m, guild, true)
	} else if strings.HasPrefix(m.Content, botPrefix+"stop") {
		go cd.voiceUsecase.StopVoice()
		cd.discord.SendMessageToChannel(m.ChannelID, "Ok -_o_-")
	} else if strings.HasPrefix(m.Content, botPrefix+"play") {
		var commandArgs []string = strings.Split(m.Content, " ")
		if len(commandArgs) > 1 {
			cd.musicUsecase.PlayYoutubeURL(commandArgs[1], s, m, guild)
		}
	} else {
		cd.discord.SendMessageToChannel(m.ChannelID, botPrefix+"help to see commands")
	}
}
