package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ssouthcity/failsafe/newsfeed/serializers"
	"golang.org/x/exp/slog"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		slog.Error("invalid amqp connection string", slog.Any("err", err))
		return
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		slog.Error("unable to connect to rabbitmq", slog.Any("err", err))
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare("failsafe.discord", false, false, false, false, nil)
	if err != nil {
		slog.Error("unable to declare queue", slog.Any("err", err))
		return
	}

	err = channel.QueueBind("failsafe.discord", "", "failsafe.newsfeed", false, nil)
	if err != nil {
		slog.Error("unable to bind queue", slog.Any("err", err))
		return
	}

	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		slog.Error("unable to consume newsfeed channel", slog.Any("err", err))
		return
	}

	serializer, err := serializers.NewFromContentType("application/msgpack")
	if err != nil {
		slog.Error("unable to get correct serializer", slog.Any("err", err))
		return
	}

	for msg := range msgs {
		story, err := serializer.Decode(msg.Body)
		if err != nil {
			slog.Error("unable to deserialize story", slog.Any("err", err))
			return
		}

		slog.Info("new story", slog.Group("story",
			slog.String("title", story.Article.Headline),
			slog.String("url", story.Article.Url),
			slog.Time("published", story.Article.PublishedAt),
		))
	}
}
