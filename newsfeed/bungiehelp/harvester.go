package bungiehelp

import (
	"context"
	"fmt"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/ssouthcity/failsafe/newsfeed"
)

const TWITTER_ROOT_URL = "https://twitter.com"

type BungieHelpHarvester struct {
	scraper          *twitterscraper.Scraper
	accountName      string
	pollingFrequency time.Duration
}

func NewHarvester(pollingFrequency time.Duration) *BungieHelpHarvester {
	return &BungieHelpHarvester{
		scraper:          twitterscraper.New(),
		accountName:      "BungieHelp",
		pollingFrequency: pollingFrequency,
	}
}

func (harvester *BungieHelpHarvester) HarvestNews(ctx context.Context, out chan newsfeed.Story) {
	ticker := time.NewTicker(harvester.pollingFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			tweets := harvester.scraper.GetTweets(context.Background(), harvester.accountName, 10)

			for tweet := range tweets {
				headline := fmt.Sprintf("Tweet from %s @ %s", harvester.accountName, tweet.TimeParsed.Format("15:04:05 01/02/06"))

				tweetUrl := fmt.Sprintf("%s/%s/status/%s", TWITTER_ROOT_URL, harvester.accountName, tweet.ID)

				thumbnail := ""
				if len(tweet.Photos) > 0 {
					thumbnail = tweet.Photos[0].URL
				}

				story := newsfeed.Story{
					Article: newsfeed.Article{
						Headline:    headline,
						Content:     tweet.Text,
						Url:         tweetUrl,
						Thumbnail:   thumbnail,
						PublishedAt: tweet.TimeParsed,
					},
					Source: newsfeed.Source{
						Name: harvester.accountName,
						Url:  fmt.Sprintf("%s/%s", TWITTER_ROOT_URL, harvester.accountName),
					},
				}

				out <- story
			}
		}
	}
}
