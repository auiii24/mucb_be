package record

type GroupRecordAnswer struct {
	QuestionGroup string `json:"questionGroup" binding:"required"`
	Score         int    `json:"score" binding:"required"`
	QuestionSize  int    `json:"questionSize" binding:"required"`
}

type CardRecordAnswer struct {
	Card string `json:"card" binding:"required"`
}

type CreateGroupRecordRequest struct {
	GroupCode *string             `json:"groupCode,omitempty" binding:"omitempty,max=64"`
	Answers   []GroupRecordAnswer `json:"answers" binding:"required,dive"`
}

type CreateManyCardRequest struct {
	GroupCode *string            `json:"groupCode,omitempty" binding:"omitempty,max=64"`
	Answers   []CardRecordAnswer `json:"answers" binding:"required,min=1,max=5,dive,required"`
}

type CreateStoryRequest struct {
	GroupCode *string `json:"groupCode,omitempty" binding:"omitempty,max=64"`
	Content   string  `json:"content" binding:"required,max=4048"`
}
