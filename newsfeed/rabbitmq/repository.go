package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/serializers"
)

type RabbitStoryRepository struct {
	channel      *amqp.Channel
	exchangeName string
	contentType  string
}

func NewStoryRepository(rabbitAddr string, exchangeName string, contentType string) (*RabbitStoryRepository, error) {
	conn, err := amqp.Dial(rabbitAddr)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(exchangeName, amqp.ExchangeFanout, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &RabbitStoryRepository{channel, exchangeName, contentType}, nil
}

func (repository *RabbitStoryRepository) SaveStory(story newsfeed.Story) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	serializer, err := serializers.NewFromContentType(repository.contentType)
	if err != nil {
		return err
	}

	body, err := serializer.Encode(story)
	if err != nil {
		return err
	}

	err = repository.channel.PublishWithContext(ctx,
		repository.exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: repository.contentType,
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
