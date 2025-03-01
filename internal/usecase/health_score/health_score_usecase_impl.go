package health_score

import (
	"mucb_be/internal/domain/health_score"
	"mucb_be/internal/domain/image"
	"mucb_be/internal/errors"
	"net/http"
)

type HealthScoreUseCaseImpl struct {
	healthScoreRepo health_score.HealthScoreRepository
	imageRepo       image.ImageRepository
}

func NewHealthScoreUseCase(
	healthScoreRepo health_score.HealthScoreRepository,
	imageRepo image.ImageRepository,
) HealthScoreInterface {
	return &HealthScoreUseCaseImpl{
		healthScoreRepo: healthScoreRepo,
		imageRepo:       imageRepo,
	}
}

func (u *HealthScoreUseCaseImpl) CreateHealthScore(req *CreateHealthScore) error {
	newHealthScore := health_score.NewHealthScore(req.Contents, req.MaximumPercent)
	err := u.healthScoreRepo.CreateHealthScore(newHealthScore)
	if err != nil {
		if err.Error() == "duplicate maximum percent" {
			return errors.NewCustomError(
				http.StatusBadRequest,
				"UCE008001001",
				"Maximum percent duplicated.",
				err.Error(),
			)
		}

		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE008001002",
			"Failed to insert health score.",
			err.Error(),
		)
	}

	return nil
}

func (u *HealthScoreUseCaseImpl) FindAllHealthScore() (*GetAllHealthScoreOutout, error) {
	healthScores, err := u.healthScoreRepo.FindAllHealthScore()
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE008002001",
			"Failed to get health score.",
			err.Error(),
		)
	}

	output := GetAllHealthScoreOutout{
		Items: healthScores,
	}

	return &output, nil
}

func (u *HealthScoreUseCaseImpl) FindHealthScoreById(id string) (*GetHealthScoreByIdOutout, error) {
	healthScoreExist, err := u.healthScoreRepo.FindHealthScoreById(id)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE008003001",
			"Failed to get health score.",
			err.Error(),
		)
	}

	return &GetHealthScoreByIdOutout{
		HealthScore: healthScoreExist,
	}, nil
}

func (u *HealthScoreUseCaseImpl) UpdateHealthScoreById(req *UpdateHealthScoreByIdRequest) error {
	healthScoreExist, err := u.healthScoreRepo.FindHealthScoreById(req.HealthScore)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE008004001",
			"Failed to get health score.",
			err.Error(),
		)
	}

	err = u.healthScoreRepo.UpdateHealthScoreById(req.HealthScore, req.Contents, req.MaximumPercent)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE008004002",
			"Failed to update health score.",
			err.Error(),
		)
	}

	for _, content := range healthScoreExist.Contents {
		if content.ContentType == health_score.ContentTypeImage {
			_ = u.imageRepo.UpdateImageStatusById(content.Content, false)
		}
	}

	for _, content := range req.Contents {
		if content.ContentType == health_score.ContentTypeImage {
			_ = u.imageRepo.UpdateImageStatusById(content.Content, true)
		}
	}

	return nil
}

func (u *HealthScoreUseCaseImpl) FindContentByScore(req *GetContentByScoreRequest) (*GetContentByScoreOutput, error) {
	healthScore, err := u.healthScoreRepo.FindContentByScore(req.Score)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE008004002",
			"Failed to update health score.",
			err.Error(),
		)
	}

	return &GetContentByScoreOutput{
		HealthScore: healthScore,
	}, nil
}
