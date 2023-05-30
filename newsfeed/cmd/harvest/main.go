package main

import (
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/mock"
	"github.com/ssouthcity/failsafe/newsfeed/terminal"
)

func main() {
	mockSrc := mock.NewHarvester(newsfeed.Source{
		Name: "Mock Fast",
	}, time.Second*5)

	mockSrc2 := mock.NewHarvester(newsfeed.Source{
		Name: "Mock Slow",
	}, time.Second*10)

	aggregatorSrc := newsfeed.NewAggregator().
		AddSource(mockSrc).
		AddSource(mockSrc2)

	terminalRepo := terminal.NewRepository()

	collector := newsfeed.NewDispatcher(aggregatorSrc, terminalRepo)

	collector.ListenForNews()
}
