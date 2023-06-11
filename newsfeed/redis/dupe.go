package redis

import (
	"context"
	"crypto/sha256"

	"github.com/go-redis/redis/v8"
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
)

type DupeStore struct {
	client *redis.Client
}

func NewDupeStore(addr string) *DupeStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &DupeStore{client}
}

func (s *DupeStore) storyHash(story newsfeed.Story) string {
	hashContent := story.Article.URL + story.Source.URL
	hashedBytes := sha256.New().Sum([]byte(hashContent))
	return string(hashedBytes)
}

func (s *DupeStore) Register(story newsfeed.Story) error {
	hash := s.storyHash(story)
	return s.client.SAdd(context.Background(), "storyhashes", hash).Err()
}

func (s *DupeStore) HasSeen(story newsfeed.Story) bool {
	hash := s.storyHash(story)
	return s.client.SIsMember(context.Background(), "storyhashes", hash).Val()
}
