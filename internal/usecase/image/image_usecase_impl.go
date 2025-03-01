package image

import (
	"mucb_be/internal/domain/image"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
)

type ImageUseCaseImpl struct {
	imageRepo image.ImageRepository
}

func NewImageUseCase(
	imageRepo image.ImageRepository,
) ImageInterface {
	return &ImageUseCaseImpl{
		imageRepo: imageRepo,
	}
}

func (u *ImageUseCaseImpl) CreateImage(req UploadImageRequest, claims *security.AccessTokenModel) (*UploadImageOutput, error) {
	newImage := image.NewImage(req.Path, req.Name, req.OriginalName, req.ContentType, req.Width, req.Height)

	err := u.imageRepo.CreateImage(newImage)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE006001001",
			"Failed to insert image.",
			err.Error(),
		)
	}

	output := UploadImageOutput{
		Image: *newImage,
	}

	return &output, nil
}

func (u *ImageUseCaseImpl) FindImageById(id string) (*FindImageOutput, error) {
	imageExist, err := u.imageRepo.FindImageByID(id)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE006002001",
			"Failed to find image.",
			err.Error(),
		)
	}

	output := FindImageOutput{
		Image: *imageExist,
	}

	return &output, nil
}
