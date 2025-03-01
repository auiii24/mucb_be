package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase admin.AdminUseCase
}

func NewAdminHandler(adminUseCase admin.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUseCase}
}

func (h AdminHandler) CreateAdmin(c *gin.Context) {
	var request admin.CreateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.adminUseCase.CreateAdmin(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
