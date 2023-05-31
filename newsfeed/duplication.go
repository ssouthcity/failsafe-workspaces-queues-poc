package newsfeed

import (
	"context"

	"golang.org/x/exp/slog"
)

type DuplicateStoryRepository interface {
	AddStory(Story) error
	HasStory(Story) bool
}

type DuplicateStoryFilter struct {
	store DuplicateStoryRepository
	input NewsHarvester
}

func NewDuplicateFilter(store DuplicateStoryRepository, input NewsHarvester) *DuplicateStoryFilter {
	return &DuplicateStoryFilter{store, input}
}

func (filter *DuplicateStoryFilter) HarvestNews(ctx context.Context, output chan Story) {
	filterChannel := make(chan Story)

	go filter.input.HarvestNews(ctx, filterChannel)

	for {
		select {
		case <-ctx.Done():
			return
		case story := <-filterChannel:
			if filter.store.HasStory(story) {
				continue
			}

			err := filter.store.AddStory(story)
			if err != nil {
				slog.Error("unable to add story to dupe store", slog.Any("err", err))
				continue
			}

			output <- story
		}
	}
}
