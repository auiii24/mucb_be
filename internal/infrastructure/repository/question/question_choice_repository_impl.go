package repository

import (
	"context"
	"fmt"
	"mucb_be/internal/domain/question"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestionChoiceRepositoryMongo struct {
	questionChoiceCollection *mongo.Collection
}

func NewQuestionChoiceRepositoryMongo(questionChoiceCollection *mongo.Collection) question.QuestionChoiceRepository {
	return &QuestionChoiceRepositoryMongo{
		questionChoiceCollection: questionChoiceCollection,
	}
}

func (r *QuestionChoiceRepositoryMongo) CreateQuestionChoice(questionChoice *question.QuestionChoice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.questionChoiceCollection.InsertOne(ctx, questionChoice)
	return err
}

func (r *QuestionChoiceRepositoryMongo) FindAllQuestionChoiceByQuestionGroup(questionGroup *primitive.ObjectID) (*[]question.QuestionChoice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.questionChoiceCollection.Find(ctx, bson.M{
		"question_group": questionGroup,
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	questions := make([]question.QuestionChoice, 0)
	err = cursor.All(ctx, &questions)
	if err != nil {
		return nil, err
	}

	return &questions, nil
}

func (r *QuestionChoiceRepositoryMongo) UpdateQuestionChoiceById(id, question string, shouldInvert bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"question":      question,
			"should_invert": shouldInvert,
			"updated_at":    time.Now(),
		},
	}

	result, err := r.questionChoiceCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("data not found")
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("can not update")
	}

	return nil
}

func (r *QuestionChoiceRepositoryMongo) RemoveChoiceById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}

	_, err = r.questionChoiceCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionChoiceRepositoryMongo) RemoveChoicesByQuestionGroupId(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"question_group": objectID,
	}

	_, err = r.questionChoiceCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
