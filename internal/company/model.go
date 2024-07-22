package company

import (
	"github.com/google/uuid"
	"github.com/mogw/micro-company/internal/kafka"
)

type Company struct {
	ID                uuid.UUID `bson:"_id" json:"id"`
	Name              string    `bson:"name" json:"name"`
	Description       string    `bson:"description,omitempty" json:"description,omitempty"`
	AmountOfEmployees int       `bson:"amount_of_employees" json:"amount_of_employees"`
	Registered        bool      `bson:"registered" json:"registered"`
	Type              string    `bson:"type" json:"type"`
}

func (c *Company) ToKafkaCompany() kafka.Company {
	return kafka.Company{
		ID:                c.ID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		Type:              c.Type,
	}
}
