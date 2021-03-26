package entities

type Feed struct {
	URL      string `json:"url"`
	Provider string `json:"provider"`
	Category string `json:"category"`
	Enabled  bool   `json:"-"`
}

type Feeds []Feed
