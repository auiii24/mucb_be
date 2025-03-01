package user

import (
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
)

type UserUseCaseImpl struct {
	userRepo   user.UserRepository
	jwtService security.JwtServiceInterface
}

func NewUserUseCase(
	userRepo user.UserRepository,
	jwtService security.JwtServiceInterface,
) UserUseCaseInterface {
	return &UserUseCaseImpl{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *UserUseCaseImpl) UpdateUserInfo(req *UpdateUserInfoRequest, claims *security.AccessTokenModel) (*UpdateUserInfoOutput, error) {
	if claims.Role != user.RoleUser {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE003001001",
			"Internal server error.",
			"",
		)
	}

	existUser, err := u.userRepo.FindUserById(claims.ID)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE003001002",
			"User not found.",
			err.Error(),
		)
	}

	err = u.userRepo.UpdateUserInfo(existUser.ID.Hex(), req.Name, req.GroupCode)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE003001004",
			"Can not update user info.",
			err.Error(),
		)
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		existUser.ID.Hex(),
		user.RoleUser,
	)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE003001003",
			"Internal server error",
			err.Error(),
		)
	}

	return &UpdateUserInfoOutput{
		AccessToken: accessToken,
	}, nil
}
