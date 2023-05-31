package inmem

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/ssouthcity/failsafe/newsfeed"
)

type InmemDuplicateStoryStore struct {
	hashes map[string]struct{}
}

func NewDupeStore() *InmemDuplicateStoryStore {
	return &InmemDuplicateStoryStore{make(map[string]struct{})}
}

func (store *InmemDuplicateStoryStore) computeHash(story newsfeed.Story) string {
	hashValue := story.Source.Name + story.Article.Headline
	hash := md5.Sum([]byte(hashValue))
	return hex.EncodeToString(hash[:])
}

func (store *InmemDuplicateStoryStore) AddStory(story newsfeed.Story) error {
	hash := store.computeHash(story)
	store.hashes[hash] = struct{}{}
	return nil
}

func (store *InmemDuplicateStoryStore) HasStory(story newsfeed.Story) bool {
	hash := store.computeHash(story)
	_, ok := store.hashes[hash]
	return ok
}
