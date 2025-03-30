package record

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRecordRepository interface {
	CreateManyGroupRecord(questionGroup *[]GroupRecord) error
	HasSubmittedToday(user primitive.ObjectID) (bool, error)
	RemoveDataByUserId(id string) error
}
