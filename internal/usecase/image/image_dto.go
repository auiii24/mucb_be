package image

import "mucb_be/internal/domain/image"

type UploadImageRequest struct {
	Path         string
	Name         string
	OriginalName string
	Width        int
	Height       int
	ContentType  string
}

type UploadImageOutput struct {
	image.Image
}

type FindImageOutput struct {
	image.Image
}
