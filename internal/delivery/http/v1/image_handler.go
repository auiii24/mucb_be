package v1

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mucb_be/internal/errors"
	imageUseCase "mucb_be/internal/usecase/image"
	"mucb_be/internal/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	imageUseCase imageUseCase.ImageInterface
}

func NewImageHandler(imageUseCase imageUseCase.ImageInterface) *ImageHandler {
	return &ImageHandler{
		imageUseCase: imageUseCase,
	}
}

func (h ImageHandler) UploadImage(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", "File type not allowed", "File type not allowed"),
		)
		return
	}

	fileOpen, err := file.Open()
	if err != nil {
		c.Error(
			errors.NewCustomError(http.StatusInternalServerError, "VE001002", "Failed to open file", err.Error()),
		)
		return
	}
	defer fileOpen.Close()

	img, _, err := image.DecodeConfig(fileOpen)
	if err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001003", "Invalid image format", err.Error()),
		)
		return
	}

	ext = filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	savePath := filepath.Join("uploads", uniqueFilename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	request := imageUseCase.UploadImageRequest{
		Path:         savePath,
		Name:         uniqueFilename,
		OriginalName: file.Filename,
		Width:        img.Width,
		Height:       img.Height,
		ContentType:  file.Header.Get("Content-Type"),
	}

	response, err := h.imageUseCase.CreateImage(request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h ImageHandler) GetImage(c *gin.Context) {
	imageID := c.Param("imageId")

	response, err := h.imageUseCase.FindImageById(imageID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Content-Type", response.ContentType)
	c.Header("X-Image-Width", fmt.Sprintf("%d", response.Width))
	c.Header("X-Image-Height", fmt.Sprintf("%d", response.Height))

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}

	c.File(response.Path)
}
