package delivery

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"discordbot-golang/internal/discord"
	jackettUsecase "discordbot-golang/internal/jackett/usecase"
	musicUsecase "discordbot-golang/internal/music/usecase"
	voiceUsecase "discordbot-golang/internal/voice/usecase"
)

//Delivery interface
type Delivery interface {
	GetCommandsHandler(*discordgo.Session, *discordgo.MessageCreate)
}

type commandsDelivery struct {
	musicUsecase   musicUsecase.Usecase
	voiceUsecase   voiceUsecase.Usecase
	jackettUsecase jackettUsecase.Usecase
	discord        discord.Discord
}

//NewCommandsDelivery new message delivery
func NewCommandsDelivery(discord discord.Discord, mu musicUsecase.Usecase, vu voiceUsecase.Usecase, ju jackettUsecase.Usecase) Delivery {
	return &commandsDelivery{
		musicUsecase:   mu,
		voiceUsecase:   vu,
		jackettUsecase: ju,
		discord:        discord,
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

	contentArray := strings.Split(m.Content, " ")
	command := strings.TrimLeft(contentArray[0], botPrefix)
	commandArgs := ""
	if len(contentArray) > 1 {
		commandArgs = strings.Join(contentArray[1:], " ")
	}
	switch command {
	case "help":
		help := []string{
			"** HELP **",
			"==============================",
			fmt.Sprintf("%splay [Youtube Link] : Plays some music", botPrefix),
			fmt.Sprintf("%sstop : Stop playing", botPrefix),
			fmt.Sprintf("%sjoin : Joins your channel", botPrefix),
			"==============================",
		}
		cd.discord.SendMessageToChannel(m.ChannelID, strings.Join(help, "\n"))
	case "join":
		cd.voiceUsecase.ConnectToVoiceChannel(s, m, guild, true)
	case "torrent":
		if commandArgs != "" {
			cd.jackettUsecase.LookupTorrent(commandArgs, s, m, guild)
		}
	case "stop":
		go cd.voiceUsecase.StopVoice()
		cd.discord.SendMessageToChannel(m.ChannelID, "Ok... -.-'")
	case "play":
		if commandArgs != "" {
			cd.musicUsecase.PlayYoutubeURL(commandArgs, s, m, guild)
		}
	default:
		cd.discord.SendMessageToChannel(m.ChannelID, fmt.Sprintf("%shelp to see available commands", botPrefix))
	}
}
