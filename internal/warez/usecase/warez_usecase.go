package usecase

import (
	"discordbot-golang/internal/discord"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	ngrok "github.com/ngrok/ngrok-api-go"
)

var (
	NGROK_CHAN chan bool = make(chan bool, 1)
)

// Usecase interface
type Usecase interface {
	CloseServer(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
	GetConsoles(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
	GetRetrocontent(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild)
	GetMovies(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
	GetTVShows(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
	GetPC(*discordgo.Session, *discordgo.MessageCreate, *discordgo.Guild)
}

type warezUsecase struct {
	discord discord.Discord
}

func NewWarezUsecase(discord discord.Discord) Usecase {
	return &warezUsecase{
		discord: discord,
	}
}

func (ju warezUsecase) GetPublicNgrokURL() (string, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:4040/api/tunnels", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	time.Sleep(5 * time.Second)
	var retries = 0
	for {
		if retries > 20 {
			return "", errors.New("Exceeded retries to API. Impossible to fetch ngrok URL.")
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			fmt.Println("API not online yet...")
			retries++
			time.Sleep(2 * time.Second)
			continue
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		var m ngrok.TunnelList
		jsonErr := json.Unmarshal(body, &m)
		if jsonErr != nil {
			return "", err
		}
		return m.Tunnels[0].PublicURL, nil
	}

}

func (ju warezUsecase) StartServer(path string, s *discordgo.Session, m *discordgo.MessageCreate) {
	cmd := exec.Command("ngrok", "http", "-auth=jack:sparrow", fmt.Sprintf("file://%s", path))

	closeCmd := func(c *exec.Cmd) {
		if err := c.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' tunnel initiated... generating link", path))
	url, err := ju.GetPublicNgrokURL()
	if err != nil {
		fmt.Println(err)
		ju.RegenerateChan(m)
		closeCmd(cmd)
		log.Fatal(err)
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Login: jack | Password: sparrow (%s)", url))
	<-NGROK_CHAN
	if err := cmd.Process.Kill(); err != nil {
		ju.RegenerateChan(m)
		log.Fatal("failed to kill process: ", err)
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' tunnel finalized.", path))
}

func (ju warezUsecase) RegenerateChan(m *discordgo.MessageCreate) {
	ju.discord.SendMessageToChannel(m.ChannelID, "Dropping any active servers...")
	if len(NGROK_CHAN) > 0 {
		go func() { NGROK_CHAN <- true }()
	}
}

func (ju warezUsecase) CloseServer(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.discord.SendMessageToChannel(m.ChannelID, "Dropping any active servers...")
	go func() { NGROK_CHAN <- true }()
}

func (ju warezUsecase) GetConsoles(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.RegenerateChan(m)
	go ju.StartServer("/mnt/beelzebub/Consoles", s, m)
}
func (ju warezUsecase) GetMovies(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.RegenerateChan(m)
	go ju.StartServer("/mnt/asmodeus/Movies", s, m)
}
func (ju warezUsecase) GetTVShows(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.RegenerateChan(m)
	go ju.StartServer("/mnt/asmodeus/TVShows", s, m)
}
func (ju warezUsecase) GetPC(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.RegenerateChan(m)
	go ju.StartServer("/mnt/beelzebub/WindowsCollection", s, m)
}
func (ju warezUsecase) GetRetrocontent(s *discordgo.Session, m *discordgo.MessageCreate, guild *discordgo.Guild) {
	ju.RegenerateChan(m)
	go ju.StartServer("/mnt/beelzebub/Retrocontent", s, m)
}
