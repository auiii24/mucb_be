package repository

import (
	"context"
	"mucb_be/internal/domain/record"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRecordRepositoryMongo struct {
	groupRecordCollection *mongo.Collection
}

func NewGroupRecordRepositoryMongo(groupRecordCollection *mongo.Collection) record.GroupRecordRepository {
	return &GroupRecordRepositoryMongo{
		groupRecordCollection: groupRecordCollection,
	}
}

func (r *GroupRecordRepositoryMongo) CreateManyGroupRecord(questionGroup *[]record.GroupRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var records []interface{}
	for _, q := range *questionGroup {
		records = append(records, q)
	}

	_, err := r.groupRecordCollection.InsertMany(ctx, records)
	return err
}

func (r *GroupRecordRepositoryMongo) HasSubmittedToday(user primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tz, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(tz)

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, tz)

	startOfDayUTC := startOfDay.UTC()
	endOfDayUTC := endOfDay.UTC()

	filter := bson.M{
		"user": user,
		"created_at": bson.M{
			"$gte": startOfDayUTC,
			"$lt":  endOfDayUTC,
		},
	}

	count, err := r.groupRecordCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
