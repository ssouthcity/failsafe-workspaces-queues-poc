package newsfeed

import "time"

type Article struct {
	Headline    string
	Content     string
	Url         string
	Thumbnail   string
	PublishedAt time.Time
}
