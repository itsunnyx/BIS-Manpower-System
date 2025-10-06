package service

import (
	"manpower/internal/domain"
	"manpower/internal/repository"
)

type ManagerService struct {
	repo repository.RequestRepo
}

func NewManagerService(repo repository.RequestRepo) *ManagerService {
	return &ManagerService{repo: repo}
}

func (s *ManagerService) ListPendingRequests() ([]domain.ManpowerRequest, error) {
	return s.repo.ListPendingForManager()
}
