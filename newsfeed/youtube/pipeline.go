package youtube

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
	"golang.org/x/exp/slog"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func FetchYoutubePlaylist(apiToken string, playlistID string, input <-chan time.Time) <-chan *youtube.PlaylistItem {
	output := make(chan *youtube.PlaylistItem)

	go func() {
		defer close(output)

		client, err := youtube.NewService(context.Background(), option.WithAPIKey(apiToken))
		if err != nil {
			slog.Error("unable to create youtube client, aborting",
				slog.Any("err", err),
			)
			return
		}

		for range input {
			call := client.PlaylistItems.List([]string{"snippet"})
			call.PlaylistId(playlistID)

			resp, err := call.Do()
			if err != nil {
				slog.Error("unable to fetch youtube videos",
					slog.String("playlist", playlistID),
					slog.Any("err", err),
				)
				continue
			}

			for i := len(resp.Items) - 1; i >= 0; i-- {
				output <- resp.Items[i]
			}
		}
	}()

	return output
}

func MapVideoToArticle(input <-chan *youtube.PlaylistItem) <-chan newsfeed.Article {
	output := make(chan newsfeed.Article)

	go func() {
		defer close(output)

		for video := range input {
			videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Snippet.ResourceId.VideoId)

			videoPublishedDate, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
			if err != nil {
				slog.Error("unable to parse youtube video time format",
					slog.String("time", video.Snippet.PublishedAt),
					slog.String("format", time.RFC3339),
					slog.Any("err", err),
				)
				continue
			}

			article := newsfeed.Article{
				ID:          uuid.New(),
				Headline:    video.Snippet.Title,
				Content:     video.Snippet.Description,
				Image:       video.Snippet.Thumbnails.Default.Url,
				URL:         videoUrl,
				PublishedAt: videoPublishedDate,
			}

			output <- article
		}
	}()

	return output
}
