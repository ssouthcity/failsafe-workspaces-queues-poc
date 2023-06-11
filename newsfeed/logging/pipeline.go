package logging

import (
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
	"golang.org/x/exp/slog"
)

func LogStoryMiddleware(input <-chan newsfeed.Story) <-chan newsfeed.Story {
	output := make(chan newsfeed.Story)

	go func() {
		defer close(output)

		for story := range input {
			slog.Info("new story",
				slog.String("headline", story.Article.Headline),
				slog.String("source", story.Source.Name),
				slog.String("url", story.Article.URL),
				slog.Any("category", story.Category),
			)

			output <- story
		}
	}()

	return output
}
