package terminal

import (
	"github.com/ssouthcity/failsafe/newsfeed"
	"golang.org/x/exp/slog"
)

type TerminalRepository struct{}

func NewRepository() *TerminalRepository {
	return &TerminalRepository{}
}

func (r *TerminalRepository) SaveStory(story newsfeed.Story) {
	slog.Info("story received",
		slog.String("headline", story.Article.Headline),
		slog.String("source", story.Source.Name),
		slog.Time("time", story.Article.PublishedAt),
	)
}
