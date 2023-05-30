package newsfeed

import "context"

type NewsAggregator struct {
	harvesters []NewsHarvester
}

func NewAggregator() *NewsAggregator {
	return &NewsAggregator{}
}

func (a *NewsAggregator) AddSource(harvester NewsHarvester) *NewsAggregator {
	a.harvesters = append(a.harvesters, harvester)
	return a
}

func (a *NewsAggregator) HarvestNews(ctx context.Context, out chan Story) {
	for _, producer := range a.harvesters {
		go producer.HarvestNews(ctx, out)
	}
}
