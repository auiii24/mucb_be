package repository

import (
	"context"
	"errors"
	"mucb_be/internal/domain/card"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CardRepositoryMongo struct {
	cardCollection *mongo.Collection
}

func NewCardRepositoryMongo(cardCollection *mongo.Collection) card.CardRepository {
	return &CardRepositoryMongo{
		cardCollection: cardCollection,
	}
}

func (r *CardRepositoryMongo) CreateCard(card *card.Card) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.cardCollection.InsertOne(ctx, card)
	return err
}

func (r *CardRepositoryMongo) FindCardById(id string) (*card.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result card.Card
	err = r.cardCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *CardRepositoryMongo) FindAllCardByRole(page, limit int, isAdmin bool) (*[]card.Card, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset := (page - 1) * limit

	filter := bson.M{}
	if !isAdmin {
		filter["is_active"] = true
	}

	total, err := r.cardCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.cardCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	cards := make([]card.Card, 0)
	if err := cursor.All(ctx, &cards); err != nil {
		return nil, 0, err
	}

	return &cards, int(total), nil

}

func (r *CardRepositoryMongo) FindCardByIdAndActivate(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  true,
			"updated_at": time.Now(),
		},
	}

	result := r.cardCollection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return errors.New("card not found")
		}
		return result.Err()
	}

	return nil
}

func (r *CardRepositoryMongo) UpdateCardById(id string, name string, description string, image primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"description": description,
			"image":       image,
			"updated_at":  time.Now(),
		},
	}

	result := r.cardCollection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return errors.New("card not found")
		}
		return result.Err()
	}

	return nil
}
