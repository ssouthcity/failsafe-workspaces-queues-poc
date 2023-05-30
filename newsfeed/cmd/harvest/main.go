package main

import (
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/bungieblog"
	"github.com/ssouthcity/failsafe/newsfeed/bungieblog/inmem"
	"github.com/ssouthcity/failsafe/newsfeed/mock"
	"github.com/ssouthcity/failsafe/newsfeed/terminal"
)

func main() {
	mockHarvester := mock.NewHarvester(newsfeed.Source{
		Name: "Mock Fast",
	}, time.Second*5)

	bungieBlogHarvester := bungieblog.NewHarvester(10*time.Second, inmem.NewStore())

	aggregator := newsfeed.NewAggregator().
		AddSource(mockHarvester).
		AddSource(bungieBlogHarvester)

	terminalRepo := terminal.NewRepository()

	dispatcher := newsfeed.NewDispatcher(aggregator, terminalRepo)

	dispatcher.ListenForNews()
}
