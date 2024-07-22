package company

import (
	"context"

	"github.com/google/uuid"
	"github.com/mogw/micro-company/internal/kafka"
)

type Service interface {
	CreateCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, id uuid.UUID, update map[string]interface{}) error
	DeleteCompany(ctx context.Context, id uuid.UUID) error
	GetCompany(ctx context.Context, id uuid.UUID) (*Company, error)
}

type service struct {
	repo     Repository
	producer kafka.Producer
}

func NewService(repo Repository, producer kafka.Producer) Service {
	return &service{repo, producer}
}

func (s *service) CreateCompany(ctx context.Context, company *Company) error {
	company.ID = uuid.New()
	if err := s.repo.CreateCompany(ctx, company); err != nil {
		return err
	}

	event := kafka.CompanyEvent{
		Type:    "create",
		Company: company.ToKafkaCompany(),
	}

	return s.produceEvent(event)
}

func (s *service) UpdateCompany(ctx context.Context, id uuid.UUID, update map[string]interface{}) error {
	if err := s.repo.UpdateCompany(ctx, id, update); err != nil {
		return err
	}

	company, err := s.repo.GetCompany(ctx, id)
	if err != nil {
		return err
	}

	event := kafka.CompanyEvent{
		Type:    "update",
		Company: company.ToKafkaCompany(),
	}

	return s.produceEvent(event)
}

func (s *service) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	company, err := s.repo.GetCompany(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteCompany(ctx, id); err != nil {
		return err
	}

	event := kafka.CompanyEvent{
		Type:    "delete",
		Company: company.ToKafkaCompany(),
	}

	return s.produceEvent(event)
}

func (s *service) GetCompany(ctx context.Context, id uuid.UUID) (*Company, error) {
	return s.repo.GetCompany(ctx, id)
}

func (s *service) produceEvent(event kafka.CompanyEvent) error {
	return nil
	// eventBytes, err := json.Marshal(event)
	// if err != nil {
	// 	return err
	// }

	// return s.producer.Produce("company-events", event.Company.ID[:], eventBytes)
}
