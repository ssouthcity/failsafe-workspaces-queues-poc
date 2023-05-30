package inmem

import "github.com/ssouthcity/failsafe/newsfeed/bungieblog"

type InmemDuplicatePostStore struct {
	ids map[string]struct{}
}

func NewStore() *InmemDuplicatePostStore {
	return &InmemDuplicatePostStore{make(map[string]struct{})}
}

func (store *InmemDuplicatePostStore) AddPost(post bungieblog.RssPost) error {
	store.ids[post.Id] = struct{}{}
	return nil
}

func (store *InmemDuplicatePostStore) HasSeenPost(post bungieblog.RssPost) bool {
	_, ok := store.ids[post.Id]
	return ok
}
