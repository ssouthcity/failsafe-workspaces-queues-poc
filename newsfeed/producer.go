package newsfeed

import "context"

type NewsHarvester interface {
	HarvestNews(context.Context, chan Story)
}
