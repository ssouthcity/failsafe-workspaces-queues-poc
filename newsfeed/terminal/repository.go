package terminal

import (
	"github.com/ssouthcity/failsafe/newsfeed"
	"golang.org/x/exp/slog"
)

type TerminalRepository struct{}

func NewRepository() *TerminalRepository {
	return &TerminalRepository{}
}

func (r *TerminalRepository) SaveArticle(article *newsfeed.Article) {
	slog.With("article", article.Title).Info("new article")
}
