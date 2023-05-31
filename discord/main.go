package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		slog.Error("invalid amqp connection string", slog.String("error", err.Error()))
		return
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		slog.Error("unable to connect to rabbitmq", slog.String("error", err.Error()))
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare("failsafe.discord", false, false, false, false, nil)
	if err != nil {
		slog.Error("unable to declare queue", slog.String("error", err.Error()))
		return
	}

	err = channel.QueueBind("failsafe.discord", "", "failsafe.newsfeed", false, nil)
	if err != nil {
		slog.Error("unable to bind queue", slog.String("error", err.Error()))
		return
	}

	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		slog.Error("unable to consume newsfeed channel", slog.String("error", err.Error()))
		return
	}

	for msg := range msgs {
		stringified := string(msg.Body)
		slog.Info(stringified)
	}
}
