package health_score

import "mucb_be/internal/domain/health_score"

type CreateHealthScore struct {
	Contents       []health_score.HealthScoreContent `json:"contents" binding:"required,dive,required"`
	MaximumPercent int                               `json:"maximumPercent" binding:"required,min=1,max=100"`
}

type GetAllHealthScoreOutout struct {
	Items *[]health_score.HealthScore `json:"items"`
}

type GetHealthScoreByIdOutout struct {
	*health_score.HealthScore
}

type UpdateHealthScoreByIdRequest struct {
	HealthScore    string                            `json:"healthScore" binding:"required"`
	Contents       []health_score.HealthScoreContent `json:"contents" binding:"required,dive,required"`
	MaximumPercent int                               `json:"maximumPercent" binding:"required,min=1,max=100"`
}

type GetContentByScoreRequest struct {
	Score int `json:"score" binding:"required,min=0,max=100"`
}

type GetContentByScoreOutput struct {
	*health_score.HealthScore
}
