package newsfeed

type NewsRepository interface {
	SaveArticle(*Article)
}
