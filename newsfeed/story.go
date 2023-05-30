package newsfeed

import "time"

type Story struct {
	Article     Article
	Source      Source
	PublishedAt time.Time
}
