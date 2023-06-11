package newsfeed

import (
	"sync"
)

func RemoveSeenStories(dupeStore DupeStoryStore, input <-chan Story) <-chan Story {
	output := make(chan Story)

	go func() {
		defer close(output)

		for story := range input {
			if dupeStore.HasSeen(story) {
				continue
			}

			output <- story
			dupeStore.Register(story)
		}
	}()

	return output
}

func MergeStreams(inputs ...<-chan Story) <-chan Story {
	output := make(chan Story, len(inputs))
	wg := sync.WaitGroup{}

	for _, input := range inputs {
		wg.Add(1)

		go func(input <-chan Story) {
			defer wg.Done()

			for story := range input {
				output <- story
			}
		}(input)
	}

	go func() {
		defer close(output)
		wg.Wait()
	}()

	return output
}

func MapArticleToStory(source Source, category Category, input <-chan Article) <-chan Story {
	output := make(chan Story)

	go func() {
		defer close(output)

		for article := range input {
			output <- Story{
				Article:  article,
				Source:   source,
				Category: category,
			}
		}
	}()

	return output
}
