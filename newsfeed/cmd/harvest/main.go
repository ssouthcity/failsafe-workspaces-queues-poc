package main

import (
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/bungieblog"
	"github.com/ssouthcity/failsafe/newsfeed/bungiehelp"
	"github.com/ssouthcity/failsafe/newsfeed/inmem"
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

	bungieHelpHarvester := bungiehelp.NewHarvester(15 * time.Second)

	bungieBlogHarvester := bungieblog.NewHarvester(15 * time.Second)

	aggregator := newsfeed.NewAggregator().
		AddSource(bungieHelpHarvester).
		AddSource(bungieBlogHarvester)

	dupeFilter := newsfeed.NewDuplicateFilter(inmemDupeStore, aggregator)

	dispatcher := newsfeed.NewDispatcher(dupeFilter, rabbitRepo)

	dispatcher.ListenForNews()
}
