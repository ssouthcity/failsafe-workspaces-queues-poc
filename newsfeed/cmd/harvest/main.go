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
	mockSrc := mock.NewHarvester(newsfeed.Source{
		Name: "Mock Fast",
	}, time.Second*5)

	bungieBlogSrc := bungieblog.NewHarvester(10*time.Second, inmem.NewStore())

	aggregatorSrc := newsfeed.NewAggregator().
		AddSource(mockSrc).
		AddSource(bungieBlogSrc)

	terminalRepo := terminal.NewRepository()

	collector := newsfeed.NewDispatcher(aggregatorSrc, terminalRepo)

	collector.ListenForNews()
}
