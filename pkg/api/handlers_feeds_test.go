package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

			rawURL := baseURL

			queryParams := []string{}
			if test.Enabled != nil {
				queryParams = append(queryParams, fmt.Sprintf("enabled=%t", *test.Enabled))
			}
			if test.Provider != "" {
				queryParams = append(queryParams, fmt.Sprintf("provider=%s", test.Provider))
			}
			if test.Category != "" {
				queryParams = append(queryParams, fmt.Sprintf("category=%s", test.Category))
			}

			if len(queryParams) != 0 {
				for i, qp := range queryParams {
					if i == 0 {
						rawURL += "?"
					} else {
						rawURL += "&"
					}
					rawURL += qp
				}
			}

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

	mockFeedsFn := func(provider string, category string, enabled bool) (feeds entities.Feeds) {
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

	// Error condition
	call := mockDB.On("GetFeeds", "errorCond", "errorCond", true)
	call = call.Return(nil, &repository.DBServiceError{})

	// For every other case
	call = call.On("GetFeeds", mock.Anything, mock.Anything, mock.Anything)
	call = call.Return(mockFeedsFn, nil)

	return mockDB
}
