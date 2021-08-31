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
	warezUsecase "discordbot-golang/internal/warez/usecase"
)

//Delivery interface
type Delivery interface {
	GetCommandsHandler(*discordgo.Session, *discordgo.MessageCreate)
}

type commandsDelivery struct {
	musicUsecase   musicUsecase.Usecase
	voiceUsecase   voiceUsecase.Usecase
	jackettUsecase jackettUsecase.Usecase
	warezUsecase   warezUsecase.Usecase
	discord        discord.Discord
}

//NewCommandsDelivery new message delivery
func NewCommandsDelivery(discord discord.Discord, mu musicUsecase.Usecase, vu voiceUsecase.Usecase, ju jackettUsecase.Usecase, wu warezUsecase.Usecase) Delivery {
	return &commandsDelivery{
		musicUsecase:   mu,
		voiceUsecase:   vu,
		jackettUsecase: ju,
		warezUsecase:   wu,
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
			fmt.Sprintf("%sleave : Leaves your channel", botPrefix),
			fmt.Sprintf("%swarez-close : Closes any open server", botPrefix),
			fmt.Sprintf("%swarez-consoles : Starts a fileserver of roms", botPrefix),
			fmt.Sprintf("%swarez-retro : Starts a fileserver of old magazines and books", botPrefix),
			fmt.Sprintf("%swarez-pc : Starts a fileserver of pc games", botPrefix),
			fmt.Sprintf("%swarez-movies : Starts a fileserver of movies", botPrefix),
			fmt.Sprintf("%swarez-tvshows : Starts a fileserver of tv shows", botPrefix),
			fmt.Sprintf("%storrent : Lookup for torrents :pirate_flag:", botPrefix),
			"==============================",
		}
		cd.discord.SendMessageToChannel(m.ChannelID, strings.Join(help, "\n"))
	case "join":
		go cd.voiceUsecase.ConnectToVoiceChannel(s, m, guild, true)
	case "torrent":
		if commandArgs != "" {
			go cd.jackettUsecase.LookupTorrent(commandArgs, s, m, guild)
		}
	case "warez-close":
		go cd.warezUsecase.CloseServer(s, m, guild)
	case "warez-retro":
		go cd.warezUsecase.GetRetrocontent(s, m, guild)
	case "warez-consoles":
		go cd.warezUsecase.GetConsoles(s, m, guild)
	case "warez-movies":
		go cd.warezUsecase.GetMovies(s, m, guild)
	case "warez-tvshows":
		go cd.warezUsecase.GetTVShows(s, m, guild)
	case "warez-pc":
		go cd.warezUsecase.GetPC(s, m, guild)
	case "leave":
		go s.ChannelVoiceJoinManual(guild.ID, "", false, true)
	case "stop":
		go cd.voiceUsecase.StopVoice()
		go s.ChannelVoiceJoinManual(guild.ID, "", false, true)
		cd.discord.SendMessageToChannel(m.ChannelID, "Ok... -.-'")
	case "play":
		if commandArgs != "" {
			go cd.musicUsecase.PlayYoutubeURL(commandArgs, s, m, guild)
		}
	default:
		cd.discord.SendMessageToChannel(m.ChannelID, fmt.Sprintf("%shelp to see available commands", botPrefix))
	}
}
