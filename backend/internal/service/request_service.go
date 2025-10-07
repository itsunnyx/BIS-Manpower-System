package service

import (
    "manpower/internal/domain"
    "manpower/internal/repository"
)

type RequesterService struct {
    repo repository.RequestRepo
}

func NewRequesterService(repo repository.RequestRepo) *RequesterService {
    return &RequesterService{repo: repo}
}

func (s *RequesterService) CreateRequest(req *domain.ManpowerRequest) error {
    return s.repo.Create(req)
}

func (s *RequesterService) GetRequestByID(id int) (*domain.ManpowerRequest, error) {
    return s.repo.GetByID(id)
}

func (s *RequesterService) ListRequests() ([]domain.ManpowerRequest, error) {
    return s.repo.List()
}
func (s *RequesterService) ListPendingRequests() ([]domain.ManpowerRequest, error) {
	return s.repo.ListPendingForManager()
}