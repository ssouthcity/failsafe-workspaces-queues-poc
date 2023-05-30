package bungieblog

type DuplicatePostStore interface {
	AddPost(RssPost) error
	HasSeenPost(RssPost) bool
}
