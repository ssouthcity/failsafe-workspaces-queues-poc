package newsfeed

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID
	Headline    string
	Content     string
	Image       string
	URL         string
	PublishedAt time.Time
}
