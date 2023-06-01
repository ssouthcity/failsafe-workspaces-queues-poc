package newsfeed

type Story struct {
	Article Article
	Source  Source
}

type StoryRepository interface {
	SaveStory(Story) error
}

type StorySerializer interface {
	Encode(Story) ([]byte, error)
	Decode([]byte) (Story, error)
}
