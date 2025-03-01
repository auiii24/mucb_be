package health_score

type HealthScoreRepository interface {
	CreateHealthScore(healthScore *HealthScore) error
	FindAllHealthScore() (*[]HealthScore, error)
	FindHealthScoreById(id string) (*HealthScore, error)
	UpdateHealthScoreById(id string, contents []HealthScoreContent, maximumPercent int) error
	FindContentByScore(score int) (*HealthScore, error)
}
