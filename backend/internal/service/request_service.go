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

func (s *RequestService) UpdateHRStatus(id int, status string) error {
	return s.repo.UpdateHRStatus(id, status)
}

func (s *RequestService) UpdateManagerStatus(id int, status string) error {
	return s.repo.UpdateManagerStatus(id, status)
}

func (s *RequestService) UpdateApproverStatus(id int, status string) error {
	return s.repo.UpdateApproverStatus(id, status)
}

func (s *RequestService) DeleteRequest(id int) error {
	return s.repo.Delete(id)
}
