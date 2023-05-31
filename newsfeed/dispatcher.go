package newsfeed

import "context"

type Dispatcher struct {
	harvester  NewsHarvester
	repository StoryRepository
}

func NewDispatcher(harvester NewsHarvester, repository StoryRepository) *Dispatcher {
	return &Dispatcher{harvester, repository}
}

func (dispatcher *Dispatcher) ListenForNews() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storyChannel := make(chan Story)

	go dispatcher.harvester.HarvestNews(ctx, storyChannel)

	for {
		select {
		case <-ctx.Done():
			return
		case story := <-storyChannel:
			dispatcher.repository.SaveStory(story)
		}
	}
}
