package image

type ImageRepository interface {
	CreateImage(image *Image) error
	FindImageByID(id string) (*Image, error)
	UpdateImageStatusById(id string, currentStats bool) error
}
