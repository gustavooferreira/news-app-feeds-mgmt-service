package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
)

func (s *Server) GetFeeds(c *gin.Context) {
	queryParams := struct {
		Enabled bool `form:"enabled"`
	}{
		Enabled: true,
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing query parameters: %s", err.Error()))
		RespondWithError(c, 400, err.Error())
		return
	}

	// Call DB ------------------------
	fq := core.FeedQuery{Provider: "", Category: "", Enabled: queryParams.Enabled}
	feeds, err := s.Repo.GetFeeds(fq)
	if err != nil {
		RespondWithError(c, 500, err.Error())
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

	// New to make sure URL is absolute and scheme is HTTP or HTTPS
	if !core.IsValideAbsoluteURL(bodyData.URL) {
		s.Logger.Info("url provided is not valid")
		RespondWithError(c, 404, "url provided is not valid")
		return
	}

	feed := core.Feed{
		URL:      bodyData.URL,
		Provider: bodyData.Provider,
		Category: bodyData.Category,
		Enabled:  true,
	}

	// Call DB ------------------------
	s.Logger.Info(fmt.Sprintf("Feed: %+v", feed))

	err = s.Repo.AddFeed(feed)
	if err != nil {
		RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func (s *Server) DeleteFeed(c *gin.Context) {
	// Validate this is a URL
	url := c.Param("url")
	url = strings.TrimPrefix(url, "/")

	// New to make sure URL is absolute and scheme is HTTP or HTTPS
	if !core.IsValideAbsoluteURL(url) {
		s.Logger.Info("url provided is not valid")
		RespondWithError(c, 404, "url provided is not valid")
		return
	}

	// Call DB ------------------------
	s.Logger.Info(fmt.Sprintf("URL : %s", url))

	c.JSON(200, gin.H{"status": "success"})
}

func (s *Server) SetFeedState(c *gin.Context) {
	url := c.Param("url")
	url = strings.TrimPrefix(url, "/")

	// New to make sure URL is absolute and scheme is HTTP or HTTPS
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

	// Call DB ------------------------
	s.Logger.Info(fmt.Sprintf("URL : %s - Body: %+v", url, *bodyData.Enabled))

	c.JSON(200, gin.H{"status": "success"})
}
