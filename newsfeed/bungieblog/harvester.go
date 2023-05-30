package bungieblog

import (
	"context"
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
)

type BungieBlogHarvester struct {
	pollingRate        time.Duration
	duplicatePostStore DuplicatePostStore
}

func NewHarvester(pollingRate time.Duration, duplicatePostStore DuplicatePostStore) *BungieBlogHarvester {
	return &BungieBlogHarvester{pollingRate, duplicatePostStore}
}

func (harvester *BungieBlogHarvester) HarvestNews(ctx context.Context, out chan newsfeed.Story) {
	ticker := time.NewTicker(harvester.pollingRate)
	defer ticker.Stop()

	rssFeed := createBungieBlogStream(ctx, harvester.pollingRate)
	filteredRssFeed := removeSeenPosts(harvester.duplicatePostStore, rssFeed)
	articleFeed := mapPostToArticle(filteredRssFeed)

	for {
		select {
		case <-ctx.Done():
			return
		case article := <-articleFeed:
			story := newsfeed.Story{
				Article: article,
				Source: newsfeed.Source{
					Name: "Bungie Blog",
					Url:  BUNGIE_RSS_ENDPOINT,
				},
			}

			out <- story
		}
	}
}
