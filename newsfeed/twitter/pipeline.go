package twitter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
	"golang.org/x/exp/slog"
)

func FetchTweets(user string, count int, input <-chan time.Time) <-chan *twitterscraper.Tweet {
	output := make(chan *twitterscraper.Tweet)

	go func() {
		defer close(output)

		scraper := twitterscraper.New()

		for range input {
			tweets, _, err := scraper.FetchTweets(user, count, "")
			if err != nil {
				slog.Warn("unable to fetch tweets",
					slog.String("user", user),
					slog.Int("maxTweets", count),
					slog.Any("err", err),
				)
				continue
			}

			for i := len(tweets) - 1; i >= 0; i-- {
				output <- tweets[i]
			}
		}
	}()

	return output
}

func MapTweetToArticle(input <-chan *twitterscraper.Tweet) <-chan newsfeed.Article {
	output := make(chan newsfeed.Article)

	go func() {
		defer close(output)

		for tweet := range input {
			output <- newsfeed.Article{
				ID:          uuid.New(),
				Headline:    fmt.Sprintf("Tweet from %s @ %d:%d", tweet.Username, tweet.TimeParsed.Hour(), tweet.TimeParsed.Minute()),
				Content:     tweet.Text,
				URL:         tweet.PermanentURL,
				PublishedAt: tweet.TimeParsed,
			}
		}
	}()

	return output
}
