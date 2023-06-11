package amqp

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
	"golang.org/x/exp/slog"
)

func PlaceStoriesInExchange(addr string, exchange string, input <-chan newsfeed.Story) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		slog.Error("unable to connect to amqp connection, aborting",
			slog.String("address", addr),
			slog.Any("err", err),
		)
		return
	}

	channel, err := conn.Channel()
	if err != nil {
		slog.Error("unable to open amqp channel, aborting",
			slog.Any("err", err),
		)
		return
	}

	err = channel.ExchangeDeclare(exchange, amqp.ExchangeFanout, false, false, false, false, nil)
	if err != nil {
		slog.Error("unable to declare channel, aborting",
			slog.String("exchange", exchange),
			slog.Any("err", err),
		)
		return
	}

	for story := range input {
		stringified, err := json.Marshal(story)
		if err != nil {
			slog.Error("unable to serialize story",
				slog.String("story", story.Article.Headline),
				slog.String("source", story.Source.Name),
				slog.Any("err", err),
			)
			continue
		}

		err = channel.PublishWithContext(context.Background(), exchange, "newsstory", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        stringified,
		})
		if err != nil {
			slog.Error("unable to publish ")
		}
	}
}
