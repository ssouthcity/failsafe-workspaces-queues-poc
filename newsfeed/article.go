package newsfeed

import "time"

type Article struct {
	Headline    string
	PublishedAt time.Time
}
