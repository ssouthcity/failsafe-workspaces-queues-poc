package rss

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
	"golang.org/x/exp/slog"
)

func FetchRssFeed(endpoint string, input <-chan time.Time) <-chan RssResponse {
	output := make(chan RssResponse)

	go func() {
		defer close(output)

		for range input {
			response, err := http.Get(endpoint)
			if err != nil {
				slog.Warn("unable to fetch rss endpoint",
					slog.String("endpoint", endpoint),
					slog.Any("err", err),
				)
				continue
			}

			rssBody := RssResponse{}

			decoder := xml.NewDecoder(response.Body)
			decoder.DefaultSpace = "RssDefault"

			err = decoder.Decode(&rssBody)
			if err != nil {
				slog.Warn("invalid rss response body",
					slog.String("endpoint", endpoint),
					slog.Any("err", err),
				)
			}

			output <- rssBody
		}
	}()

	return output
}

func MapRssFeedToArticles(input <-chan RssResponse) <-chan newsfeed.Article {
	output := make(chan newsfeed.Article)

	go func() {
		defer close(output)

		for rssFeed := range input {

			domain, err := url.Parse(rssFeed.Channel.Link)
			if err != nil {
				slog.Error("unable to parse rss domain",
					slog.String("rssfeed", rssFeed.Channel.Title),
					slog.Any("err", err),
				)
				continue
			}

			for i := len(rssFeed.Channel.Items) - 1; i >= 0; i-- {
				post := rssFeed.Channel.Items[i]

				pubdate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", post.PubDate)
				if err != nil {
					slog.Error("invalid time format",
						slog.String("rssfeed", rssFeed.Channel.Title),
						slog.String("article", post.Title),
						slog.Any("err", err),
					)
					continue
				}

				output <- newsfeed.Article{
					ID:          uuid.New(),
					Headline:    post.Title,
					Content:     post.Description,
					Image:       post.Image,
					URL:         fmt.Sprintf("https://%s%s", domain.Host, post.Link),
					PublishedAt: pubdate,
				}
			}
		}
	}()

	return output
}
