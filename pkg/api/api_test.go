package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/api"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/log"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetFeeds(t *testing.T) {

	logger := log.NullLogger{}
	db := repository.NewDatabase()

	server := api.NewServer("", 9999, false, logger, db)
	router := server.Router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/feeds", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{}`, w.Body.String())
}
