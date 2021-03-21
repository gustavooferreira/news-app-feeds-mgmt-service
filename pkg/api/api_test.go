package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/mocks"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/api"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetFeeds(t *testing.T) {
	assert := assert.New(t)

	logger := log.NullLogger{}

	data := dbData()
	mockDB := setupMockDB(data)

	expectedFeeds, err := json.Marshal(data)
	require.NoError(t, err)

	server := api.NewServer("", 9999, false, logger, mockDB)
	router := server.Router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/feeds", nil)
	router.ServeHTTP(w, req)

	assert.Equal(200, w.Code)
	assert.JSONEq(string(expectedFeeds), w.Body.String())
}

func dbData() core.Feeds {
	data := core.Feeds{
		core.Feed{
			URL:      "url1",
			Provider: "provider1",
			Category: "category1",
			Enabled:  true},
	}

	return data
}

func setupMockDB(data core.Feeds) *mocks.Repository {
	mockDB := &mocks.Repository{}

	mockResultFn := func(fq core.FeedQuery) core.Feeds {
		feeds := data
		return feeds
	}

	mockDB.On("GetFeeds", mock.Anything).Return(mockResultFn, nil)

	return mockDB
}
