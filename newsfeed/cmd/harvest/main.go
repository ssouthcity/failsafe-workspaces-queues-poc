package main

import (
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/bungieblog"
	"github.com/ssouthcity/failsafe/newsfeed/inmem"
	"github.com/ssouthcity/failsafe/newsfeed/mock"
	"github.com/ssouthcity/failsafe/newsfeed/rabbitmq"
	"golang.org/x/exp/slog"
)

func main() {
	inmemDupeStore := inmem.NewDupeStore()

	rabbitRepo, err := rabbitmq.NewStoryRepository("amqp://guest:guest@localhost:5672", "failsafe.newsfeed")
	if err != nil {
		slog.Error("unable to connect to rabbitmq", slog.Any("err", err))
		return
	}

	mockHarvester := mock.NewHarvester(newsfeed.Source{
		Name: "Second Counter",
	}, time.Second*1)

	bungieBlogHarvester := bungieblog.NewHarvester(10 * time.Second)

	aggregator := newsfeed.NewAggregator().
		AddSource(mockHarvester).
		AddSource(bungieBlogHarvester)

	dupeFilter := newsfeed.NewDuplicateFilter(inmemDupeStore, aggregator)

	dispatcher := newsfeed.NewDispatcher(dupeFilter, rabbitRepo)

	dispatcher.ListenForNews()
}
