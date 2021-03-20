package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
)

func (s *Server) GetFeeds(c *gin.Context) {
	// enabled := c.Param("enabled")
	// enabled := c.DefaultQuery("enabled", "true")

	fq := core.FeedQuery{Provider: "", Category: "", Enabled: false}
	feeds, err := s.Repo.GetFeeds(fq)
	if err != nil {
		RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, feeds)
}

func (s *Server) AddFeed(c *gin.Context) {

	var feed core.Feed

	err := c.ShouldBindJSON(&feed)
	if err != nil {
		// return error
		RespondWithError(c, 400, err.Error())
		return
	}

	// fmt.Printf("Query: %+v", feed)

	// Validate input

	err = s.Repo.AddFeed(feed)
	if err != nil {
		RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func (s *Server) DeleteFeed(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

func (s *Server) SetFeedState(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}
