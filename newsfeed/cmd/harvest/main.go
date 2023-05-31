package main

import (
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/bungieblog"
	"github.com/ssouthcity/failsafe/newsfeed/inmem"
	"github.com/ssouthcity/failsafe/newsfeed/mock"
	"github.com/ssouthcity/failsafe/newsfeed/terminal"
)

func main() {
	inmemDupeStore := inmem.NewDupeStore()

	terminalRepo := terminal.NewRepository()

	mockHarvester := mock.NewHarvester(newsfeed.Source{
		Name: "Second Counter",
	}, time.Second*1)

	bungieBlogHarvester := bungieblog.NewHarvester(10 * time.Second)

	aggregator := newsfeed.NewAggregator().
		AddSource(mockHarvester).
		AddSource(bungieBlogHarvester)

	dupeFilter := newsfeed.NewDuplicateFilter(inmemDupeStore, aggregator)

	dispatcher := newsfeed.NewDispatcher(dupeFilter, terminalRepo)

	dispatcher.ListenForNews()
}
