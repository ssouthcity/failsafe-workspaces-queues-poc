package bungieblog

import (
	"encoding/xml"
	"net/http"
)

const BUNGIE_RSS_ENDPOINT = "https://www.bungie.net/en/rss/News"

type RssResponse struct {
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Items []RssPost `xml:"item"`
}

type RssPost struct {
	Id          string `xml:"guid"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Image       string `xml:"image"`
}

func fetchRssFeed() (RssResponse, error) {
	rssBody := RssResponse{}

	response, err := http.Get(BUNGIE_RSS_ENDPOINT)
	if err != nil {
		return rssBody, err
	}

	err = xml.NewDecoder(response.Body).Decode(&rssBody)
	if err != nil {
		return rssBody, err
	}

	return rssBody, nil
}
