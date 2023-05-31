package newsfeed

type Story struct {
	Article Article
	Source  Source
}

type StoryRepository interface {
	SaveStory(Story) error
}
