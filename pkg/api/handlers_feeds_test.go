package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/mocks"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/api"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/entities"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/log"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetFeedsHandler(t *testing.T) {
	type ErrorResponseBody struct {
		Message string `json:"message"`
	}

	falseV := false
	trueV := true

	assert := assert.New(t)

	logger := log.NullLogger{}
	mockDB := setupMockDB()
	server := api.NewServer("", 9999, false, logger, mockDB)
	router := server.Router

	baseURL := "/api/v1/feeds"

	tests := map[string]struct {
		Enabled              *bool
		Provider             string
		Category             string
		expectedStatusCode   int
		expectedResponseBody interface{}
	}{
		"test 1": {
			Provider:           "errorCond",
			Category:           "errorCond",
			Enabled:            &trueV,
			expectedStatusCode: 500,
			expectedResponseBody: ErrorResponseBody{
				Message: "Internal error",
			}},
		"test 2": {
			Provider:           "",
			Category:           "",
			Enabled:            &trueV,
			expectedStatusCode: 200,
			expectedResponseBody: entities.Feeds{
				entities.Feed{
					URL:      "http://feeds.bbci.co.uk/news/technology/rss.xml",
					Provider: "BBC News",
					Category: "Technology"},
				entities.Feed{
					URL:      "http://feeds.bbci.co.uk/news/uk/rss.xml",
					Provider: "BBC News",
					Category: "UK"},
				entities.Feed{
					URL:      "http://feeds.skynews.com/feeds/rss/technology.xml",
					Provider: "Sky News",
					Category: "Technology"}}},
		"test 3": {
			Provider:           "Sky News",
			Category:           "",
			Enabled:            &falseV,
			expectedStatusCode: 200,
			expectedResponseBody: entities.Feeds{
				entities.Feed{
					URL:      "http://feeds.skynews.com/feeds/rss/uk.xml",
					Provider: "Sky News",
					Category: "UK",
					Enabled:  false}}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			rawURL := BuildQueryParams(baseURL, test.Provider, test.Category, test.Enabled)

			req, err := http.NewRequest("GET", rawURL, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(test.expectedStatusCode, w.Code)

			switch w.Code {
			case 200:
				responseBody := entities.Feeds{}
				err = json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				assert.Equal(test.expectedResponseBody, responseBody)
			case 500:
				responseBody := ErrorResponseBody{}
				err = json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				assert.Equal(test.expectedResponseBody, responseBody)
			}
		})
	}
}

func TestAddFeedHandler(t *testing.T) {
	assert := assert.New(t)

	logger := log.NullLogger{}
	mockDB := setupMockDB()
	server := api.NewServer("", 9999, false, logger, mockDB)
	router := server.Router

	baseURL := "/api/v1/feeds"

	tests := map[string]struct {
		URL                string `json:"url"`
		Provider           string `json:"provider"`
		Category           string `json:"category"`
		expectedStatusCode int    `json:"-"`
	}{
		"test1": {
			expectedStatusCode: 400,
		},
		"test 2": {
			URL:                "http://example.com",
			expectedStatusCode: 400,
		},
		"test 3": {
			URL:                "invalid_url",
			Provider:           "provider1",
			Category:           "category1",
			expectedStatusCode: 400,
		},
		"test 4": {
			URL:                "http://feeds.bbci.co.uk/news/technology/rss.xml",
			Provider:           "provider1",
			Category:           "category1",
			expectedStatusCode: 409,
		},
		"test 5": {
			URL:                "http://errorCond.com",
			Provider:           "errorCond",
			Category:           "errorCond",
			expectedStatusCode: 500,
		},
		"test 6": {
			URL:                "http://example.com",
			Provider:           "provider1",
			Category:           "category1",
			expectedStatusCode: 204,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			requestBody, err := json.Marshal(test)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(test.expectedStatusCode, w.Code)
		})
	}
}

func TestSetFeedStateHandler(t *testing.T) {
	assert := assert.New(t)

	falseV := false
	trueV := true

	logger := log.NullLogger{}
	mockDB := setupMockDB()
	server := api.NewServer("", 9999, false, logger, mockDB)
	router := server.Router

	baseURL := "/api/v1/feeds/"

	tests := map[string]struct {
		URL                string `json:"-"`
		Enabled            *bool  `json:"enabled"`
		expectedStatusCode int    `json:"-"`
	}{
		"test1": {
			URL:                "http://feeds.bbci.co.uk/news/technology/rss.xml",
			expectedStatusCode: 400,
		},
		"test 2": {
			URL:                "invalid_url",
			Enabled:            &falseV,
			expectedStatusCode: 400,
		},
		"test 3": {
			URL:                "http://url.does.not.exist.com",
			Enabled:            &falseV,
			expectedStatusCode: 404,
		},
		"test 4": {
			URL:                "http://errorCond.com",
			Enabled:            &trueV,
			expectedStatusCode: 500,
		},
		"test 5": {
			URL:                "http://feeds.bbci.co.uk/news/technology/rss.xml",
			Enabled:            &falseV,
			expectedStatusCode: 204,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			rawURL := baseURL + test.URL

			requestBody, err := json.Marshal(test)
			require.NoError(t, err)
			req, err := http.NewRequest("PUT", rawURL, bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(test.expectedStatusCode, w.Code)
		})
	}
}

func TestDeleteFeedHandler(t *testing.T) {
	assert := assert.New(t)

	logger := log.NullLogger{}
	mockDB := setupMockDB()
	server := api.NewServer("", 9999, false, logger, mockDB)
	router := server.Router

	baseURL := "/api/v1/feeds/"

	tests := map[string]struct {
		URL                string
		expectedStatusCode int
	}{
		"test1": {
			URL:                "invalid_url",
			expectedStatusCode: 400,
		},
		"test 2": {
			URL:                "http://url.does.not.exist.com",
			expectedStatusCode: 404,
		},
		"test 4": {
			URL:                "http://errorCond.com",
			expectedStatusCode: 500,
		},
		"test 5": {
			URL:                "http://feeds.bbci.co.uk/news/technology/rss.xml",
			expectedStatusCode: 204,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			rawURL := baseURL + test.URL

			req, err := http.NewRequest("DELETE", rawURL, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(test.expectedStatusCode, w.Code)
		})
	}
}

func BuildQueryParams(rawURL string, provider string, category string, enabled *bool) string {
	v := url.Values{}

	if provider != "" {
		v.Set("provider", provider)
	}
	if category != "" {
		v.Set("category", category)
	}
	if enabled != nil {
		v.Set("enabled", strconv.FormatBool(*enabled))
	}
	queryParams := v.Encode()

	if len(v) != 0 {
		rawURL += "?" + queryParams
	}

	return rawURL
}

func GenData() entities.Feeds {
	data := entities.Feeds{
		entities.Feed{
			URL:      "http://feeds.bbci.co.uk/news/technology/rss.xml",
			Provider: "BBC News",
			Category: "Technology",
			Enabled:  true},
		entities.Feed{
			URL:      "http://feeds.bbci.co.uk/news/uk/rss.xml",
			Provider: "BBC News",
			Category: "UK",
			Enabled:  true},
		entities.Feed{
			URL:      "http://feeds.skynews.com/feeds/rss/technology.xml",
			Provider: "Sky News",
			Category: "Technology",
			Enabled:  true},
		entities.Feed{
			URL:      "http://feeds.skynews.com/feeds/rss/uk.xml",
			Provider: "Sky News",
			Category: "UK",
			Enabled:  false},
	}

	return data
}

func setupMockDB() *mocks.Repository {
	mockDB := &mocks.Repository{}

	data := GenData()

	mockGetFeedsFn := func(provider string, category string, enabled bool) (feeds entities.Feeds) {
		feeds = entities.Feeds{}

		for _, item := range data {
			if provider != "" && provider != item.Provider {
				continue
			}

			if category != "" && category != item.Category {
				continue
			}

			if enabled != item.Enabled {
				continue
			}

			feeds = append(feeds, item)
		}

		return feeds
	}

	mockAddFeedFn := func(feed entities.Feed) (err error) {
		// Error condition
		if feed.URL == "http://errorCond.com" && feed.Provider == "errorCond" && feed.Category == "errorCond" {
			return &repository.DBServiceError{}
		}

		for _, item := range data {
			if item.URL == feed.URL {
				return &repository.DBDUPError{}
			}
		}
		return nil
	}

	mockSetFeedStateFn := func(url string, enabled bool) (err error) {
		// Error condition
		if url == "http://errorCond.com" {
			return &repository.DBServiceError{}
		}

		notFound := true
		for _, item := range data {
			if item.URL == url {
				notFound = false
			}
		}

		if notFound {
			return &repository.DBNotFoundError{}
		}
		return nil
	}

	mockDeleteFeedFn := func(url string) (err error) {
		// Error condition
		if url == "http://errorCond.com" {
			return &repository.DBServiceError{}
		}

		notFound := true
		for _, item := range data {
			if item.URL == url {
				notFound = false
			}
		}

		if notFound {
			return &repository.DBNotFoundError{}
		}
		return nil
	}

	// GetFeeds mock -------------------------------------
	// Error condition
	call := mockDB.On("GetFeeds", "errorCond", "errorCond", true)
	call = call.Return(nil, &repository.DBServiceError{})

	// For every other case
	call = call.On("GetFeeds", mock.Anything, mock.Anything, mock.Anything)
	call = call.Return(mockGetFeedsFn, nil)

	// AddFeed mock -------------------------------------
	call = call.On("AddFeed", mock.Anything)
	call = call.Return(mockAddFeedFn)

	// SetFeedState mock -------------------------------------
	call = call.On("SetFeedState", mock.Anything, mock.Anything)
	call = call.Return(mockSetFeedStateFn)

	// DeleteFeed mock -------------------------------------
	call = call.On("DeleteFeed", mock.Anything)
	call = call.Return(mockDeleteFeedFn)

	return mockDB
}
