package services

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func UploadToYouTube(
	videoPath, title, description,
	clientID, clientSecret, refreshToken string,
) (string, error) {

	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{youtube.YoutubeUploadScope},
		Endpoint:     google.Endpoint,
	}

	token := &oauth2.Token{RefreshToken: refreshToken}
	client := conf.Client(ctx, token)

	service, _ := youtube.New(client)

	file, _ := os.Open(videoPath)
	defer file.Close()

	call := service.Videos.Insert(
		[]string{"snippet", "status"},
		&youtube.Video{
			Snippet: &youtube.VideoSnippet{
				Title:       title,
				Description: description,
			},
			Status: &youtube.VideoStatus{
				PrivacyStatus: "unlisted",
			},
		},
	)

	resp, err := call.Media(file).Do()
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}