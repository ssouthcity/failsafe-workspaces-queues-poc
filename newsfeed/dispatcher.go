package newsfeed

import "context"

type Dispatcher struct {
	producer   NewsHarvester
	repository NewsRepository
}

func NewDispatcher(producer NewsHarvester, repository NewsRepository) *Dispatcher {
	return &Dispatcher{producer, repository}
}

func (dispatcher *Dispatcher) ListenForNews() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storyChannel := make(chan Story)

	go dispatcher.producer.HarvestNews(ctx, storyChannel)

	for {
		select {
		case <-ctx.Done():
			return
		case story := <-storyChannel:
			dispatcher.repository.SaveStory(story)
		}
	}
}
