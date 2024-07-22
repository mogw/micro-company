package company

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mogw/micro-company/internal/kafka"
)

type Service interface {
	CreateCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, uuid uuid.UUID, update map[string]interface{}) error
	DeleteCompany(ctx context.Context, uuid uuid.UUID) error
	GetCompany(ctx context.Context, uuid uuid.UUID) (*Company, error)
}

type service struct {
	repo     Repository
	producer kafka.Producer
}

func NewService(repo Repository, producer kafka.Producer) Service {
	return &service{repo, producer}
}

func (s *service) CreateCompany(ctx context.Context, company *Company) error {
	company.UUID = uuid.New()
	if err := s.repo.CreateCompany(ctx, company); err != nil {
		return err
	}

	event := kafka.CompanyEvent{
		Type:    "create",
		Company: company.ToKafkaCompany(),
	}

	return s.produceEvent(event)
}

func (s *service) UpdateCompany(ctx context.Context, uuid uuid.UUID, update map[string]interface{}) error {
	if err := s.repo.UpdateCompany(ctx, uuid, update); err != nil {
		return err
	}

	company, err := s.repo.GetCompany(ctx, uuid)
	if err != nil {
		return err
	}

	event := kafka.CompanyEvent{
		Type:    "update",
		Company: company.ToKafkaCompany(),
	}

	return s.produceEvent(event)
}

func (s *service) DeleteCompany(ctx context.Context, uuid uuid.UUID) error {
	company, err := s.repo.GetCompany(ctx, uuid)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteCompany(ctx, uuid); err != nil {
		return err
	}

	event := kafka.CompanyEvent{
		Type:    "delete",
		Company: company.ToKafkaCompany(),
	}

	return s.produceEvent(event)
}

func (s *service) GetCompany(ctx context.Context, uuid uuid.UUID) (*Company, error) {
	return s.repo.GetCompany(ctx, uuid)
}

func (s *service) produceEvent(event kafka.CompanyEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.producer.Produce("company-events", event.Company.UUID[:], eventBytes)
}
