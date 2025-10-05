package http

import (
	"database/sql"

	"manpower/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewServer(db *sql.DB) *gin.Engine {
	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok in server.g"})
	})

	api := r.Group("/api")
	{
		api.GET("/requests", handlers.GetManpowerRequests(db))
		api.POST("/requests", handlers.CreateManpowerRequest(db))
	}

	return r
}
