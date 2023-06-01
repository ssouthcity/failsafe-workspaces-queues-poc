package serializers

import (
	"errors"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/serializers/msgpack"
)

func NewFromContentType(contentType string) (newsfeed.StorySerializer, error) {
	switch contentType {
	case "application/msgpack":
		return &msgpack.StorySerializer{}, nil
	default:
		return nil, errors.New("invalid content type")
	}
}
