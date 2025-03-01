package record

import "go.mongodb.org/mongo-driver/bson/primitive"

type CardRecordRepository interface {
	CreateManyGroupRecord(cardRecords *[]CardRecord) error
	HasSubmittedToday(user primitive.ObjectID) (bool, error)
}
