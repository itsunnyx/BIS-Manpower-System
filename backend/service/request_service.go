package service

import (
    "manpower/domain"
    "manpower/repository"
    "time"
)

type RequestService struct {
    repo repository.RequestRepository
}

func NewRequestService(repo repository.RequestRepository) *RequestService {
    return &RequestService{repo: repo}
}

func (s *RequestService) CreateRequest(req domain.Request) (domain.Request, error) {
    req.CreatedAt = time.Now()
    err := s.repo.Save(req)
    return req, err
}

func (s *RequestService) GetRequests() ([]domain.Request, error) {
    return s.repo.FindAll()
}
