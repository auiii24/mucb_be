package repository

import (
	"context"
	"mucb_be/internal/domain/admin"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepositoryMongo struct {
	adminCollection *mongo.Collection
}

func NewAdminRepositoryMongo(adminCollection *mongo.Collection) admin.AdminRepository {
	return &AdminRepositoryMongo{
		adminCollection: adminCollection,
	}
}

func (r *AdminRepositoryMongo) CreateAdmin(admin *admin.Admin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.adminCollection.InsertOne(ctx, admin)
	return err
}

func (r *AdminRepositoryMongo) FindAdminByEmail(email string) (*admin.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result admin.Admin
	err := r.adminCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AdminRepositoryMongo) FindAdminById(id string) (*admin.Admin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result admin.Admin
	err = r.adminCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
