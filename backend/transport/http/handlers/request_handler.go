package handlers

import (
    "encoding/json" 
    "backend/domain"
    "backend/service"
    "net/http"
)

type RequestHandler struct {
    service *service.RequestService
}

func NewRequestHandler(s *service.RequestService) *RequestHandler {
    return &RequestHandler{service: s}
}

func (h *RequestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
    var body domain.Request
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    req, err := h.service.CreateRequest(body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(req)
}

func (h *RequestHandler) GetRequests(w http.ResponseWriter, r *http.Request) {
    requests, err := h.service.GetRequests()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(requests)
}
