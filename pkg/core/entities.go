package core

type Feed struct {
	URL      string
	Provider string
	Category string
	Enabled  bool
}

// Index is based on the feed URL
type Feeds []Feed

// Pick one of these to query
type FeedQuery struct {
	Provider string `json:"provider"`
	Category string `json:"category"`
	Enabled  bool   `json:"enabled"`
}
