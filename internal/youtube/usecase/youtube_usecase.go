package usecase

import (
	"fmt"
	"net/url"
	"strings"

	"discordbot-golang/internal/music/model"
	"github.com/kkdai/youtube/v2"
)

//Usecase interface
type Usecase interface {
	GetYoutubeDownloadURL(string) (*model.Song, error)
}

type youtubeUsecase struct {
	ytdlClient *youtube.Client
}

//NewYoutubeUsecase new message delivery
func NewYoutubeUsecase() Usecase {
	client := &youtube.Client{}

	return &youtubeUsecase{
		ytdlClient: client,
	}
}

//GetYoutubeDownloadURL return youtube download url
func (yu youtubeUsecase) GetYoutubeDownloadURL(link string) (*model.Song, error) {
	client := yu.ytdlClient
	videoInfo, err := client.GetVideo(link)
	if err != nil {
		return nil, err
	}
	for _, format := range videoInfo.Formats {
		if strings.Contains(format.MimeType, "opus") {
			data, err := client.GetStreamURL(videoInfo, &format)
			if err != nil {
				return nil, err
			}
			thumbUrl, err := url.Parse(videoInfo.Thumbnails[0].URL)
			if err != nil {
				return nil, err
			}
			return &model.Song{
				Title:        videoInfo.Title,
				Link:         link,
				DownloadLink: data,
				Duration:     videoInfo.Duration,
				Uploader:     videoInfo.Author,
				ThumbnailURL: thumbUrl,
			}, nil
		}
	}
	return nil, fmt.Errorf("Audio format not found")
}
