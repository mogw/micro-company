package company

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, uuid uuid.UUID, update map[string]interface{}) error
	DeleteCompany(ctx context.Context, uuid uuid.UUID) error
	GetCompany(ctx context.Context, uuid uuid.UUID) (*Company, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateCompany(ctx context.Context, company *Company) error {
	company.UUID = uuid.New()
	return s.repo.CreateCompany(ctx, company)
}

func (s *service) UpdateCompany(ctx context.Context, uuid uuid.UUID, update map[string]interface{}) error {
	return s.repo.UpdateCompany(ctx, uuid, update)
}

func (s *service) DeleteCompany(ctx context.Context, uuid uuid.UUID) error {
	return s.repo.DeleteCompany(ctx, uuid)
}

func (s *service) GetCompany(ctx context.Context, uuid uuid.UUID) (*Company, error) {
	return s.repo.GetCompany(ctx, uuid)
}
