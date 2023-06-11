package rss

type RssResponse struct {
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"RssDefault link"`
	Items       []RssPost `xml:"item"`
}

type RssPost struct {
	Id          string `xml:"guid"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Image       string `xml:"image"`
}
