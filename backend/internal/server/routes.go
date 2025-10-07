package server

import (
	"database/sql"
	"manpower/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	DB     *sql.DB
}

func NewServer(db *sql.DB) *Server {
	r := gin.Default()

	// Handlers
	requestHandler := handlers.NewRequestHandler(db)

	// Group: Manager
	managerGroup := r.Group("/api/manager")
	{
		managerGroup.GET("/requests", requestHandler.GetAllRequests)
		managerGroup.POST("/requests", requestHandler.CreateRequest)
		managerGroup.PUT("/requests/:id/approve", requestHandler.ApproveRequest)
	}

	// Group: HR
	hrGroup := r.Group("/api/hr")
	{
		hrGroup.GET("/requests", requestHandler.GetAllRequests)
		hrGroup.PUT("/requests/:id", requestHandler.UpdateRequestStatus) // เช่น HR อัปเดตสถานะ
		hrGroup.PUT("/requests/:id/approve", requestHandler.ApproveRequest)
	}

	// Group: Approver
	approverGroup := r.Group("/api/approver")
	{
		approverGroup.GET("/requests", requestHandler.GetAllRequests)
		approverGroup.PUT("/requests/:id", requestHandler.UpdateRequestStatus)
		approverGroup.PUT("/requests/:id/approve", requestHandl	er.ApproveRequest)
	}

	// Group: ฝ่ายบริหาร (?)
	recruiterGroup := r.Group("/api/recruiter")
	{
		recruiterGroup.GET("/requests", requestHandler.GetAllRequests)
		recruiterGroup.PUT("/requests/:id", requestHandler.UpdateRequestStatus)
		recruiterGroup.POST("/sync/:id", requestHandler.SyncToRecruitment)
	}

	// Group: Admin
	adminGroup := r.Group("/api/admin")
	{
		adminGroup.GET("/requests", requestHandler.GetAllRequests)
		adminGroup.DELETE("/requests/:id", requestHandler.DeleteRequest)
	}

	return &Server{
		Engine: r,
		DB:     db,
	}
}

func (s *Server) Run(addr string) error {
	return s.Engine.Run(addr)
}















// func RegisterRoutes(r *gin.Engine, reqSvc *service.RequesterService, managerSvc *service.ManagerService) {

// 	api := r.Group("/api/v1")

// 	requesterHandler := handlers.NewRequesterHandler(reqSvc)
// 	managerHandler := handlers.NewManagerHandler(managerSvc)

// 	requester := api.Group("/requester")
// 	{
// 		requester.GET("/requests", requesterHandler.GetManpowerRequests)
// 	}

	// hr := api.Group("/hr")
	// {
	// 	hr.GET("/requests", hrHandler.ListAllRequests)
	// 	hr.PUT("/requests/:id/review", hrHandler.ReviewRequest)
	// }

	// manager := api.Group("/manager")
	// {
	// 	manager.GET("/requests", managerHandler.ListPendingRequests)
	// 	// manager.PUT("/requests/:id/approve", managerHandler.ApproveRequest)
	// }

	// admin := api.Group("/admin")
	// {
	// 	admin.GET("/requests", adminHandler.GetAllRequests)
	// 	admin.GET("/users", adminHandler.ListUsers)
	// }
