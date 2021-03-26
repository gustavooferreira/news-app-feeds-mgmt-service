package repository

import (
	"errors"
	"fmt"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/entities"
	"gorm.io/gorm"
)

type DBServiceError struct {
	Msg string
	Err error
}

func (e *DBServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}
	return e.Msg
}
func (e *DBServiceError) Unwrap() error {
	return e.Err
}

type DBDUPError struct{}

func (e *DBDUPError) Error() string { return "database error: duplicate entry" }

type DBNotFoundError struct{}

func (e *DBNotFoundError) Error() string { return "database error: entry not found" }

type DatabaseService struct {
	Database *Database
}

func NewDatabaseService(host string, port int, username string, password string, dbname string) (dbs *DatabaseService, err error) {
	dbs = &DatabaseService{}
	dbs.Database, err = NewDatabase(host, port, username, password, dbname)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (dbs *DatabaseService) Close() error {
	return dbs.Database.Close()
}

func (dbs *DatabaseService) HealthCheck() error {
	return dbs.Database.HealthCheck()
}

func (dbs *DatabaseService) GetFeeds(provider string, category string, enabled bool) (feeds entities.Feeds, err error) {
	feedRecords, err := dbs.Database.FindAllFeedRecords(provider, category, enabled)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Feeds{}, nil
	} else if err != nil {
		return nil, &DBServiceError{Msg: "database error", Err: err}
	}

	feedList := make(entities.Feeds, 0, len(feedRecords))

	for _, feedRecord := range feedRecords {
		feedItem := entities.Feed{
			URL:      feedRecord.URL,
			Provider: feedRecord.Provider.Name,
			Category: feedRecord.Category.Name,
			Enabled:  *feedRecord.Enabled,
		}

		feedList = append(feedList, feedItem)
	}

	return feedList, nil
}

func (dbs *DatabaseService) AddFeed(feed entities.Feed) (err error) {
	err = dbs.Database.InsertFeedRecord(feed.URL, feed.Provider, feed.Category, feed.Enabled)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
				return &DBDUPError{}
			}
		}
		return &DBServiceError{Msg: "database error", Err: err}
	}

	return nil
}

func (dbs *DatabaseService) SetFeedState(url string, enabled bool) (err error) {
	err = dbs.Database.UpdateFeedState(url, enabled)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &DBNotFoundError{}
	} else if err != nil {
		return &DBServiceError{Msg: "database error", Err: err}
	}

	return nil
}

func (dbs *DatabaseService) DeleteFeed(url string) (err error) {
	err = dbs.Database.DeleteFeedState(url)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &DBNotFoundError{}
	} else if err != nil {
		return &DBServiceError{Msg: "database error", Err: err}
	}

	return nil
}
