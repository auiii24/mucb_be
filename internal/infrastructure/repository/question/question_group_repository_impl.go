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

type QuestionGroupRepositoryMongo struct {
	questionGroupCollection *mongo.Collection
}

func NewQuestionGroupRepositoryMongo(questionGroupCollection *mongo.Collection) question.QuestionGroupRepository {
	return &QuestionGroupRepositoryMongo{
		questionGroupCollection: questionGroupCollection,
	}
}

func (r *QuestionGroupRepositoryMongo) CreateQuestionGroup(questionGroup *question.QuestionGroup) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.questionGroupCollection.InsertOne(ctx, questionGroup)
	return err
}

func (r *QuestionGroupRepositoryMongo) FindAllQuestionGroup(page, limit int) (*[]question.QuestionGroup, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	offset := (page - 1) * limit

	total, err := r.questionGroupCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.questionGroupCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	groups := make([]question.QuestionGroup, 0)
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, 0, err
	}

	return &groups, int(total), nil
}

func (r *QuestionGroupRepositoryMongo) FindGroupsWithRandomChoices() (*[]question.GroupsWithRandomChoices, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "question_choices"},
				{Key: "let", Value: bson.D{
					{Key: "groupId", Value: "$_id"},
					{Key: "groupLimit", Value: "$limit"},
				}},
				{Key: "pipeline", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$question_group", "$$groupId"}},
						}}},
					}},
					bson.D{{Key: "$sample", Value: bson.D{
						{Key: "size", Value: 5}},
					}},
				}},
				{Key: "as", Value: "choices"},
			}},
		},
		{
			{Key: "$set", Value: bson.D{
				{Key: "choices", Value: bson.D{
					{Key: "$slice", Value: bson.A{"$choices", "$limit"}}},
				},
			}},
		},
	}

	cursor, err := r.questionGroupCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var questionGroups []question.GroupsWithRandomChoices
	if err := cursor.All(ctx, &questionGroups); err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	return &questionGroups, nil
}

func (r *QuestionGroupRepositoryMongo) RemoveQuestionGroupById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}

	_, err = r.questionGroupCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionGroupRepositoryMongo) FindQuestionGroupById(id string) (*question.QuestionGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var questionGroup question.QuestionGroup
	err = r.questionGroupCollection.FindOne(ctx, bson.M{
		"_id": objectID,
	}).Decode(&questionGroup)
	if err != nil {
		return nil, err
	}

	return &questionGroup, nil
}

func (r *QuestionGroupRepositoryMongo) UpdateQuestionGroupById(id string, columnName string, description string, limit int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"column_name": columnName,
			"description": description,
			"limit":       limit,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.questionGroupCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
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
