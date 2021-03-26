package core

import (
	"context"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/entities"
)

type Repository interface {
	HealthCheck() error
	GetFeeds(provider string, category string, enabled bool) (feeds entities.Feeds, err error)
	AddFeed(feed entities.Feed) (err error)
	SetFeedState(url string, enabled bool) (err error)
	DeleteFeed(url string) (err error)
}

// ShutDowner represents anything that can be shutdown like an HTTP server.
type ShutDowner interface {
	ShutDown(ctx context.Context) error
}
