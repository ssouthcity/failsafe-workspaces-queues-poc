package main

import (
	"context"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"golang.org/x/exp/slog"
)

func main() {
	twitter := twitterscraper.New()

	tweets := twitter.GetTweets(context.Background(), "BungieHelp", 10)

	for tweet := range tweets {
		slog.Info(tweet.Text)
	}
}
