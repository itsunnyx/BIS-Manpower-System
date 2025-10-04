package http

import (
	"backend/repository"
	"backend/service"
	"backend/transport/http/handlers"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func NewServer(db *sql.DB) *http.Server {
	repo := repository.NewRequestRepository(db)
	service := service.NewRequestService(repo)
	handler := handlers.NewRequestHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/requests", handler.CreateRequest).Methods("POST")
	r.HandleFunc("/requests", handler.GetRequests).Methods("GET")

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
