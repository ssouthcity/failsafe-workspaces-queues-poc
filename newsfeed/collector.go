package newsfeed

import "context"

type Collector struct {
	source     NewsSource
	repository NewsRepository
}

func NewCollector(source NewsSource, repository NewsRepository) *Collector {
	return &Collector{source, repository}
}

func (c *Collector) ListenForNews() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	articleChan := make(chan *Article)

	go c.source.CollectNews(ctx, articleChan)

	for {
		select {
		case <-ctx.Done():
			return
		case article := <-articleChan:
			c.repository.SaveArticle(article)
		}
	}
}
