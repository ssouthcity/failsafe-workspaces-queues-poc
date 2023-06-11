package main

import (
	"flag"
	"time"

	"github.com/ssouthcity/failsafe/newsfeed/amqp"
	"github.com/ssouthcity/failsafe/newsfeed/logging"
	"github.com/ssouthcity/failsafe/newsfeed/newsfeed"
	"github.com/ssouthcity/failsafe/newsfeed/redis"
	"github.com/ssouthcity/failsafe/newsfeed/rss"
	"github.com/ssouthcity/failsafe/newsfeed/twitter"
	"github.com/ssouthcity/failsafe/newsfeed/youtube"
)

const BUNGIE_RSS_ENDPOINT = "https://www.bungie.net/en/rss/News"

const DESTINY2_YOUTUBE_ENGLISH_PLAYLIST_ID = "PLw2gyMFmq40pL-jC1jFPreWHGuV7g_Kmu"

var (
	youtubeApiKey = flag.String("youtube", "", "api key for the youtube api")
)

func main() {
	flag.Parse()

	// bungieblog
	bungieBlogSource := newsfeed.Source{
		Name: "Bungie RSS Feed",
		URL:  BUNGIE_RSS_ENDPOINT,
	}

	bungieBlogFetchTimer := time.NewTicker(5 * time.Second)
	bungieBlogRssFeed := rss.FetchRssFeed(BUNGIE_RSS_ENDPOINT, bungieBlogFetchTimer.C)
	allBungieBlogArticles := rss.MapRssFeedToArticles(bungieBlogRssFeed)
	bungieBlogStories := newsfeed.MapArticleToStory(bungieBlogSource, newsfeed.News, allBungieBlogArticles)

	// bungiehelp
	bungieNewsSource := newsfeed.Source{
		Name: "BungieHelp Twitter",
		URL:  "http://twitter.com/BungieHelp",
	}

	bungieHelpTimer := time.NewTicker(5 * time.Second)
	bungieHelpTweets := twitter.FetchTweets("BungieHelp", 50, bungieHelpTimer.C)
	bungieHelpArticles := twitter.MapTweetToArticle(bungieHelpTweets)
	bungieHelpStories := newsfeed.MapArticleToStory(bungieNewsSource, newsfeed.ServerUpdates, bungieHelpArticles)

	// destiny2team
	d2teamSource := newsfeed.Source{
		Name: "Destiny 2 Team",
		URL:  "https://twitter.com/Destiny2Team",
	}

	d2teamTimer := time.NewTicker(5 * time.Second)
	d2teamTweets := twitter.FetchTweets("Destiny2Team", 50, d2teamTimer.C)
	d2teamArticles := twitter.MapTweetToArticle(d2teamTweets)
	d2TeamStories := newsfeed.MapArticleToStory(d2teamSource, newsfeed.SocialMedia, d2teamArticles)

	// destiny 2 youtube
	d2youtubeSource := newsfeed.Source{
		Name: "Destiny 2 Youtube Channel",
		URL:  "https://www.youtube.com/@destinythegame",
	}

	d2YoutubeTimer := time.NewTicker(5 * time.Second)
	d2YoutubeVideos := youtube.FetchYoutubePlaylist(*youtubeApiKey, DESTINY2_YOUTUBE_ENGLISH_PLAYLIST_ID, d2YoutubeTimer.C)
	d2YoutubeArticles := youtube.MapVideoToArticle(d2YoutubeVideos)
	d2youtubeStories := newsfeed.MapArticleToStory(d2youtubeSource, newsfeed.Videos, d2YoutubeArticles)

	// merged sources
	mergedStories := newsfeed.MergeStreams(bungieBlogStories, bungieHelpStories, d2TeamStories, d2youtubeStories)
	newStories := newsfeed.RemoveSeenStories(redis.NewDupeStore("localhost:6379"), mergedStories)
	loggedStories := logging.LogStoryMiddleware(newStories)

	amqp.PlaceStoriesInExchange("amqp://guest:guest@localhost:5672", "failsafe", loggedStories)
}
