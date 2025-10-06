package httpx

import (

	"manpower/internal/transport/http/handlers"
	"manpower/internal/service"

	"github.com/gin-gonic/gin"

)

func RegisterRoutes(r *gin.Engine, reqSvc *service.RequestService) {

	api := r.Group("/api/v1")

	requesterHandler := handlers.NewRequestHandler(reqSvc)

	requester := api.Group("/requester")
	{
		requester.GET("/requests", requesterHandler.GetManpowerRequests)
	}

	// hr := api.Group("/hr")
	// {
	// 	hr.GET("/requests", hrHandler.ListAllRequests)
	// 	hr.PUT("/requests/:id/review", hrHandler.ReviewRequest)
	// }

	manager := api.Group("/manager")
	{
		manager.GET("/requests", requesterHandler.ListPendingRequests)
		// manager.PUT("/requests/:id/approve", managerHandler.ApproveRequest)
	}

	// admin := api.Group("/admin")
	// {
	// 	admin.GET("/requests", adminHandler.GetAllRequests)
	// 	admin.GET("/users", adminHandler.ListUsers)
	// }
}