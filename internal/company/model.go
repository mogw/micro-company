package company

import (
	"github.com/google/uuid"
	"github.com/mogw/micro-company/internal/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UUID              uuid.UUID          `bson:"uuid"`
	Name              string             `bson:"name"`
	Description       string             `bson:"description,omitempty"`
	AmountOfEmployees int                `bson:"amount_of_employees"`
	Registered        bool               `bson:"registered"`
	Type              string             `bson:"type"`
}

func (c *Company) ToKafkaCompany() kafka.Company {
	return kafka.Company{
		UUID:              c.UUID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		Type:              c.Type,
	}
}
