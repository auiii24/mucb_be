package image

import "mucb_be/internal/infrastructure/security"

type ImageInterface interface {
	CreateImage(req UploadImageRequest, claims *security.AccessTokenModel) (*UploadImageOutput, error)
	FindImageById(id string) (*FindImageOutput, error)
}
