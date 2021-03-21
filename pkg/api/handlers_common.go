package api

import (
	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "no route found"})
}

func RespondWithError(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, gin.H{"message": message})
}

// Healthcheck checks health of the service.
func (s *Server) Healthcheck(c *gin.Context) {

	// TODO: ping db to make sure it's still alive!

	c.JSON(200, gin.H{
		"status": "OK",
	})
}
