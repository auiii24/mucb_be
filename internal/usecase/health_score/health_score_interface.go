package health_score

type HealthScoreInterface interface {
	CreateHealthScore(req *CreateHealthScore) error
	FindAllHealthScore() (*GetAllHealthScoreOutout, error)
	FindHealthScoreById(id string) (*GetHealthScoreByIdOutout, error)
	UpdateHealthScoreById(req *UpdateHealthScoreByIdRequest) error
	FindContentByScore(req *GetContentByScoreRequest) (*GetContentByScoreOutput, error)
}
