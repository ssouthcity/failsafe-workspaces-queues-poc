package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
)

type MockSource struct {
	name             string
	pollingFrequency time.Duration
	iteration        int
}

func NewSource(name string, pollingFrequency time.Duration) *MockSource {
	return &MockSource{name, pollingFrequency, 0}
}

func (a *MockSource) CollectNews(ctx context.Context, out chan *newsfeed.Article) {
	for {
		a.iteration++

		article := &newsfeed.Article{
			Title: a.name + " " + strconv.Itoa(a.iteration),
		}

		select {
		case <-ctx.Done():
			return
		case out <- article:
			time.Sleep(a.pollingFrequency)
		}
	}
}
