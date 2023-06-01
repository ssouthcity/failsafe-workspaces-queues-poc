package msgpack

import (
	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/vmihailenco/msgpack/v5"
)

type StorySerializer struct{}

func (s *StorySerializer) Encode(story newsfeed.Story) ([]byte, error) {
	return msgpack.Marshal(story)
}

func (s *StorySerializer) Decode(rawStory []byte) (newsfeed.Story, error) {
	story := newsfeed.Story{}

	err := msgpack.Unmarshal(rawStory, &story)
	if err != nil {
		return story, err
	}

	return story, nil
}
