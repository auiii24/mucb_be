package question

type QuestionGroupRepository interface {
	CreateQuestionGroup(questionGroup *QuestionGroup) error
	FindAllQuestionGroup(page, limit int) (*[]QuestionGroup, int, error)
	FindGroupsWithRandomChoices() (*[]GroupsWithRandomChoices, error)
}
