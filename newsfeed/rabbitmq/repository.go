package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ssouthcity/failsafe/newsfeed"
)

type RabbitStoryRepository struct {
	channel      *amqp.Channel
	exchangeName string
}

func NewStoryRepository(rabbitAddr string, exchangeName string) (*RabbitStoryRepository, error) {
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

	return &RabbitStoryRepository{channel, exchangeName}, nil
}

func (repository *RabbitStoryRepository) SaveStory(story newsfeed.Story) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	body, err := json.Marshal(story)
	if err != nil {
		return err
	}

	err = repository.channel.PublishWithContext(ctx,
		repository.exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
