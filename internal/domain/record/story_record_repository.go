package record

type StoryRecordRepository interface {
	CreateStoryRecord(storyRecord *StoryRecord) error
}
