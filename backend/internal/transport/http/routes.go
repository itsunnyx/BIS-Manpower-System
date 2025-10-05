package httpx

import (

	"manpower/internal/transport/http/handlers"
	"manpower/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, reqSvc *service.RequestService) {

	api := r.Group("/api/v1")

	requester := api.Group("/requester")
	{
		requester.GET("/requests", handlers.GetManpowerRequests(reqSvc))
	}
}