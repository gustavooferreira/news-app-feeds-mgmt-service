package core

import "context"

type Repository interface {
	GetFeeds(fq FeedQuery) (feeds Feeds, err error)
	AddFeed(feed Feed) (err error)
	SetFeedState(enabled bool) (err error)
	DeleteFeed(url string) (err error)
}

// ShutDowner represents anything that can be shutdown like an HTTP server.
type ShutDowner interface {
	ShutDown(ctx context.Context) error
}
