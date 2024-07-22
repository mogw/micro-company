package company

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateCompany(ctx context.Context, company *Company) error
	UpdateCompany(ctx context.Context, uuid uuid.UUID, update bson.M) error
	DeleteCompany(ctx context.Context, uuid uuid.UUID) error
	GetCompany(ctx context.Context, uuid uuid.UUID) (*Company, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Client, dbName, collName string) Repository {
	collection := db.Database(dbName).Collection(collName)
	return &repository{collection}
}

func (r *repository) CreateCompany(ctx context.Context, company *Company) error {
	_, err := r.collection.InsertOne(ctx, company)
	return err
}

func (r *repository) UpdateCompany(ctx context.Context, uuid uuid.UUID, update bson.M) error {
	filter := bson.M{"uuid": uuid}
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update}, options.Update().SetUpsert(false))
	return err
}

func (r *repository) DeleteCompany(ctx context.Context, uuid uuid.UUID) error {
	filter := bson.M{"uuid": uuid}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *repository) GetCompany(ctx context.Context, uuid uuid.UUID) (*Company, error) {
	filter := bson.M{"uuid": uuid}
	var company Company
	err := r.collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}
