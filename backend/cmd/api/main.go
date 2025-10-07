package main

import (
	"manpower/internal/db"
	"manpower/internal/repository"
	"manpower/internal/service"
	httpx "manpower/internal/transport/http"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// เริ่มต้น DB
	conn := db.InitDB()
	defer conn.Close()

	reqRepo := repository.NewRequestRepo(conn)

	reqSvc := service.NewRequestService(reqRepo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		err := conn.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "unhealthy", "error": err})
			return
		}
		c.JSON(200, gin.H{"message": "healthy"})
	})

	httpx.RegisterRoutes(r, reqSvc)

	r.Run(":8080")
}
