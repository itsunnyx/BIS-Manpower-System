package server

import (
	"manpower/internal/service"
	"manpower/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, reqSvc *service.RequestService) {

	api := r.Group("/api/v1")

	requestHandler := handlers.NewRequestHandler(reqSvc)

	requester := api.Group("/requester")
	{
		requester.GET("/requests", requestHandler.GetManpowerRequests)
	}

	hr := api.Group("/hr")
	{
		hr.GET("/requests", requestHandler.GetManpowerRequests)
		hr.POST("/requests", requestHandler.CreateRequest)
		// hr.PUT("/requests/:id", requestHandler.UpdateRequestStatus) // เช่น HR อัปเดตสถานะ
		hr.PUT("/requests/:id/approve", requestHandler.HrApproveRequest)
	}

	manager := api.Group("/manager")
	{
		manager.GET("/requests", requestHandler.GetManpowerRequests)
		manager.POST("/requests", requestHandler.CreateRequest)
		manager.PUT("/requests/:id/approve", requestHandler.ManagerApproveRequest)
	}

	approver := api.Group("/approver")
	{
		approver.GET("/requests", requestHandler.GetManpowerRequests)
		approver.POST("/requests", requestHandler.CreateRequest)
		approver.PUT("/requests/:id/approve", requestHandler.ApproverApproveRequest)
	}
	admin := api.Group("/admin")
	{
		admin.GET("/requests", requestHandler.GetManpowerRequests)
	}

}

// http://127.0.0.1:8080/api/v1/manager/requests/1/approve

// http://127.0.0.1:8080/api/v1/manager/requests
