package record

type StoryRecordRepository interface {
	CreateStoryRecord(storyRecord *StoryRecord) error
	RemoveDataByUserId(id string) error
}
