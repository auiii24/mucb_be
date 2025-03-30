package question

type QuestionGroupRepository interface {
	CreateQuestionGroup(questionGroup *QuestionGroup) error
	FindAllQuestionGroup(page, limit int) (*[]QuestionGroup, int, error)
	FindGroupsWithRandomChoices() (*[]GroupsWithRandomChoices, error)
	RemoveQuestionGroupById(id string) error
	FindQuestionGroupById(id string) (*QuestionGroup, error)
	UpdateQuestionGroupById(id, columnName, description string, limit int) error
}
