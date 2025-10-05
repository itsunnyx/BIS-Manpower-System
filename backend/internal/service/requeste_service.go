package service

import (
    "manpower/internal/domain"
    "manpower/internal/repository"
)

type RequestService struct {
    repo *repository.RequestRepo
}

func NewRequestService(repo *repository.RequestRepo) *RequestService {
    return &RequestService{repo: repo}
}

func (s *RequestService) CreateRequest(req *domain.ManpowerRequest) error {
    return s.repo.Create(req)
}

func (s *RequestService) GetRequestByID(id int) (*domain.ManpowerRequest, error) {
    return s.repo.GetByID(id)
}

func (s *RequestService) ListRequests() ([]domain.ManpowerRequest, error) {
    return s.repo.List()
}