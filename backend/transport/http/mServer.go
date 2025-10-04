package http

import (
    "database/sql"
    "log"
    "backend/repository"
    "backend/service"
    "backend/transport/http/handlers"
    "net/http"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
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
