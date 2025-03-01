package record

import (
	"mucb_be/internal/infrastructure/security"
)

type RecordInterface interface {
	CreateManyGroupRecord(req *CreateGroupRecordRequest, claims *security.AccessTokenModel) error
	CreateManyCardRecord(req *CreateManyCardRequest, claims *security.AccessTokenModel) error
	CreateStoryRecord(req *CreateStoryRequest, claims *security.AccessTokenModel) error
}
