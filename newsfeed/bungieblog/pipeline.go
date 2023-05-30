package bungieblog

import (
	"context"
	"time"

	"github.com/ssouthcity/failsafe/newsfeed"
	"golang.org/x/exp/slog"
)

func createBungieBlogStream(ctx context.Context, pollingRate time.Duration) chan RssPost {
	postFeed := make(chan RssPost)

	go func() {
		ticker := time.NewTicker(pollingRate)
		defer ticker.Stop()

		defer close(postFeed)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rss, err := fetchRssFeed()
				if err != nil {
					slog.Error("unable to read rss feed", slog.String("error", err.Error()))
					continue
				}

				for _, item := range rss.Channel.Items {
					postFeed <- item
				}
			}
		}
	}()

	return postFeed
}

func removeSeenPosts(store DuplicatePostStore, input chan RssPost) chan RssPost {
	output := make(chan RssPost)

	go func() {
		for post := range input {
			if store.HasSeenPost(post) {
				slog.Debug("skipping seen post", slog.String("post", post.Title))
				continue
			}

			err := store.AddPost(post)
			if err != nil {
				slog.Error("error while adding post to duplication store", slog.String("error", err.Error()))
				continue
			}

			output <- post
		}
	}()

	return output
}

func mapPostToArticle(input chan RssPost) chan newsfeed.Article {
	output := make(chan newsfeed.Article)

	go func() {
		for post := range input {
			pubdate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", post.PubDate)
			if err != nil {
				slog.Error("invalid time format", slog.String("error", err.Error()))
				continue
			}

			article := newsfeed.Article{
				Headline:    post.Title,
				PublishedAt: pubdate,
			}

			output <- article
		}
	}()

	return output
}