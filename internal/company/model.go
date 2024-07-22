package company

import (
	"github.com/google/uuid"
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
