package newsfeed

type NewsRepository interface {
	SaveStory(Story) error
}
