package admin

import (
	"mucb_be/internal/domain/admin"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
	"strings"
)

type AdminUseCaseImpl struct {
	adminRepo   admin.AdminRepository
	hashService security.HashServiceInterface
}

func NewAdminUseCase(
	adminRepo admin.AdminRepository,
	hashService security.HashServiceInterface,
) AdminUseCase {
	return &AdminUseCaseImpl{
		adminRepo:   adminRepo,
		hashService: hashService,
	}
}

func (u *AdminUseCaseImpl) CreateAdmin(req *CreateAdminRequest) error {
	existingAdmin, _ := u.adminRepo.FindAdminByEmail(strings.ToLower(req.Email))
	if existingAdmin != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE001001001",
			"Email already exist.",
			"",
		)
	}

	hashedPassword, err := u.hashService.HashPassword(req.Password)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE001001002",
			"Can not create admin.",
			err.Error(),
		)
	}

	newAdmin := admin.NewAdmin(req.Name, req.Email, hashedPassword, req.Role)

	err = u.adminRepo.CreateAdmin(newAdmin)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE001001003",
			"Can not create admin.",
			err.Error(),
		)
	}

	return nil
}
