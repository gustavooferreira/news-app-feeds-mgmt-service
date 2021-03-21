package repository

import (
	"sync"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
)

type Database struct {
	mu    sync.Mutex // protects Feeds
	Feeds core.Feeds
}

// NewDatabase creates a new database manager
// We need to return a pointer cos the struct is holding a mutex and we don't want to copy mutexes
func NewDatabase() (*Database, error) {
	db := Database{Feeds: core.Feeds{}}
	return &db, nil
}

func (db *Database) GetFeeds(fq core.FeedQuery) (core.Feeds, error) {
	// for now return everything but we want to filter it
	return db.Feeds, nil
}

func (db *Database) AddFeed(feed core.Feed) error {
	// if _, ok := db.Feeds[feed.URL]; ok {
	// 	return fmt.Errorf("feed already existes")
	// }

	// db.Feeds[feed.URL] = feed
	return nil
}

func (db *Database) SetFeedState(enabled bool) (err error) {
	return nil
}

func (db *Database) DeleteFeed(url string) (err error) {
	return nil
}
