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
	UpdateCompany(ctx context.Context, id uuid.UUID, update bson.M) error
	DeleteCompany(ctx context.Context, id uuid.UUID) error
	GetCompany(ctx context.Context, id uuid.UUID) (*Company, error)
	IsNameUnique(ctx context.Context, name string) (bool, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Client, dbName, collName string) Repository {
	collection := db.Database(dbName).Collection(collName)
	return &repository{collection}
}

func (r *repository) CreateCompany(ctx context.Context, company *Company) error {
	isUnique, err := r.IsNameUnique(ctx, company.Name)
	if err != nil {
		return err
	}
	if !isUnique {
		return errors.New("company name must be unique")
	}

	_, err = r.collection.InsertOne(ctx, company)
	return err
}

func (r *repository) UpdateCompany(ctx context.Context, id uuid.UUID, update bson.M) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update}, options.Update().SetUpsert(false))
	return err
}

func (r *repository) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *repository) GetCompany(ctx context.Context, id uuid.UUID) (*Company, error) {
	filter := bson.M{"_id": id}
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

func (r *repository) IsNameUnique(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	var company Company
	err := r.collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
