package usecase

import (
	"context"
	"discordbot-golang/internal/discord"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/webtor-io/go-jackett"
)

//Usecase interface
type Usecase interface {
	LookupTorrent(string, *discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
}

type jacketUsecase struct {
	discord discord.Discord
}

//NewJackettUsecase new voice usecase
func NewJackettUsecase(discord discord.Discord) Usecase {
	return &jacketUsecase{
		discord: discord,
	}
}

//JoiAndPlayAudioFile return youtube download url
func (ju jacketUsecase) LookupTorrent(torrent string, s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.discord.SendMessageToChannel(m.ChannelID, fmt.Sprintf("Looking for the best 3 matches of '%s' on 1337x.to...", torrent))
	ctx := context.Background()
	j := jackett.NewJackett(&jackett.Settings{
		ApiURL: os.Getenv("JACKETT_API_URL"),
		ApiKey: os.Getenv("JACKETT_API_KEY"),
	})
	resp, err := j.Fetch(ctx, &jackett.FetchRequest{
		Trackers: []string{"1337x"},
		Query:    torrent,
	})
	if err != nil {
		panic(err)
	}
	threeFirst := []jackett.Result{
		resp.Results[0],
		resp.Results[1],
		resp.Results[2],
	}
	for _, r := range threeFirst {
		req, err := http.NewRequest("GET", r.Link, nil)
		if err != nil {
			panic(err)
		}
		client := new(http.Client)
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirect")
		}
		response, err := client.Do(req)
		if err != nil {
			if response.StatusCode == http.StatusFound {
				e := discordgo.MessageEmbed{
					URL:   r.Guid,
					Title: fmt.Sprintf("[%dmb] (%v) %v", r.Size/1024, r.CategoryDesc, r.Title),
				}
				if _, err := s.ChannelMessageSendComplex(m.ChannelID,
					&discordgo.MessageSend{Embed: &e},
				); err != nil {
					panic(err)
				}
				u, _ := response.Location()
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", u))
			} else {
				panic(err)
			}
		}
	}

}
