package main

import (
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/aggregator"
	"github.com/ssouthcity/failsafe/newsfeed/mock"
	"github.com/ssouthcity/failsafe/newsfeed/terminal"
)

func main() {
	mockSrc := mock.NewSource("Fast", time.Second*5)
	mockSrc2 := mock.NewSource("Slow", time.Second*10)

	aggregatorSrc := aggregator.NewAggregator().
		AddSource(mockSrc).
		AddSource(mockSrc2)

	terminalRepo := terminal.NewRepository()

	collector := newsfeed.NewCollector(aggregatorSrc, terminalRepo)

	collector.ListenForNews()
}
