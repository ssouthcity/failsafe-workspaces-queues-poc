package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
)

type MockHarvester struct {
	source           newsfeed.Source
	pollingFrequency time.Duration
}

func NewHarvester(source newsfeed.Source, pollingFrequency time.Duration) *MockHarvester {
	return &MockHarvester{source, pollingFrequency}
}

func (harvester *MockHarvester) HarvestNews(ctx context.Context, out chan newsfeed.Story) {
	ticker := time.NewTicker(harvester.pollingFrequency)
	defer ticker.Stop()

	i := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			i++

			story := newsfeed.Story{
				Source: harvester.source,
				Article: newsfeed.Article{
					Headline: strconv.Itoa(i),
				},
				PublishedAt: time.Now(),
			}

			out <- story
		}
	}
}
