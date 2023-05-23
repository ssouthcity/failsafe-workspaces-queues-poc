package aggregator

import (
	"context"

	"github.com/ssouthcity/failsafe/newsfeed"
)

type NewsAggregator struct {
	sources []newsfeed.NewsSource
}

func NewAggregator() *NewsAggregator {
	return &NewsAggregator{}
}

func (a *NewsAggregator) AddSource(source newsfeed.NewsSource) *NewsAggregator {
	a.sources = append(a.sources, source)
	return a
}

func (a *NewsAggregator) CollectNews(ctx context.Context, out chan *newsfeed.Article) {
	for _, producer := range a.sources {
		go producer.CollectNews(ctx, out)
	}
}
