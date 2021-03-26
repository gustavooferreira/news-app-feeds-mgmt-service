package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/entities"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/repository"
)

func (s *Server) GetFeeds(c *gin.Context) {
	queryParams := struct {
		Enabled  bool   `form:"enabled"`
		Provider string `form:"provider"`
		Category string `form:"category"`
	}{
		Enabled: true,
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing query parameters: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	feeds, err := s.Repo.GetFeeds(queryParams.Provider, queryParams.Category, queryParams.Enabled)
	if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.JSON(200, feeds)
}

func (s *Server) AddFeed(c *gin.Context) {
	bodyData := struct {
		URL      string `json:"url" binding:"required"`
		Provider string `json:"provider" binding:"required"`
		Category string `json:"category" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	// Need to make sure URL is absolute and scheme is HTTP or HTTPS
	if !core.IsValideAbsoluteURL(bodyData.URL) {
		s.Logger.Info("url provided is not valid")
		RespondWithError(c, 400, "url provided is not valid")
		return
	}

	feed := entities.Feed{
		URL:      bodyData.URL,
		Provider: bodyData.Provider,
		Category: bodyData.Category,
		Enabled:  true,
	}

	err = s.Repo.AddFeed(feed)
	if err, ok := err.(*repository.DBDUPError); ok {
		s.Logger.Error(err.Error())
		RespondWithError(c, 409, "RSS URL feed already exists in the database")
		return
	} else if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) SetFeedState(c *gin.Context) {
	url := c.Param("url")
	url = strings.TrimPrefix(url, "/")

	// Need to make sure URL is absolute and scheme is HTTP or HTTPS
	if !core.IsValideAbsoluteURL(url) {
		s.Logger.Info("url provided is not valid")
		RespondWithError(c, 404, "url provided is not valid")
		return
	}

	bodyData := struct {
		// Need to make Enabled a bool pointer to validate it exists
		Enabled *bool `json:"enabled" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	err = s.Repo.SetFeedState(url, *bodyData.Enabled)
	if err, ok := err.(*repository.DBNotFoundError); ok {
		s.Logger.Error(err.Error())
		RespondWithError(c, 404, "URL not found")
		return
	} else if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) DeleteFeed(c *gin.Context) {
	url := c.Param("url")
	url = strings.TrimPrefix(url, "/")

	// Need to make sure URL is absolute and scheme is HTTP or HTTPS
	if !core.IsValideAbsoluteURL(url) {
		s.Logger.Info("url provided is not valid")
		RespondWithError(c, 404, "url provided is not valid")
		return
	}

	err := s.Repo.DeleteFeed(url)
	if err, ok := err.(*repository.DBNotFoundError); ok {
		s.Logger.Error(err.Error())
		RespondWithError(c, 404, "URL not found")
		return
	} else if err != nil {
		s.Logger.Error(err.Error())
		RespondWithError(c, 500, "Internal error")
		return
	}

	c.Status(http.StatusNoContent)
}
