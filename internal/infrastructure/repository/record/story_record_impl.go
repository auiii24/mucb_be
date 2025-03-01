package repository

import (
	"context"
	"mucb_be/internal/domain/record"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type StoryRecordRepositoryMongo struct {
	storyRecordCollection *mongo.Collection
}

func NewStoryRecordRepositoryMongo(storyRecordCollection *mongo.Collection) record.StoryRecordRepository {
	return &StoryRecordRepositoryMongo{
		storyRecordCollection: storyRecordCollection,
	}
}

func (r *StoryRecordRepositoryMongo) CreateStoryRecord(storyRecord *record.StoryRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.storyRecordCollection.InsertOne(ctx, storyRecord)
	return err
}
