package newsfeed

type DupeStoryStore interface {
	Register(story Story) error
	HasSeen(story Story) bool
}
