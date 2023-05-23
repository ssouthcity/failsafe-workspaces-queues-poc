package newsfeed

import "context"

type NewsSource interface {
	CollectNews(context.Context, chan *Article)
}
