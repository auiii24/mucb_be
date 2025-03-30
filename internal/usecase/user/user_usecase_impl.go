package user

import (
	"mucb_be/internal/domain/auth"
	"mucb_be/internal/domain/record"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
)

type UserUseCaseImpl struct {
	userRepo        user.UserRepository
	groupRecordRepo record.GroupRecordRepository
	cardRecordRepo  record.CardRecordRepository
	storyRecordRepo record.StoryRecordRepository
	authRepo        auth.AuthRepository
	jwtService      security.JwtServiceInterface
}

func NewUserUseCase(
	userRepo user.UserRepository,
	groupRecordRepo record.GroupRecordRepository,
	cardRecordRepo record.CardRecordRepository,
	storyRecordRepo record.StoryRecordRepository,
	authRepo auth.AuthRepository,
	jwtService security.JwtServiceInterface,
) UserUseCaseInterface {
	return &UserUseCaseImpl{
		userRepo:        userRepo,
		groupRecordRepo: groupRecordRepo,
		cardRecordRepo:  cardRecordRepo,
		storyRecordRepo: storyRecordRepo,
		authRepo:        authRepo,
		jwtService:      jwtService,
	}
}

func (u *UserUseCaseImpl) UpdateUserInfo(req *UpdateUserInfoRequest, claims *security.AccessTokenModel) (*UpdateUserInfoOutput, error) {
	if claims.Role != user.RoleUser {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE003001001",
			"Failed to check role.",
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

func (u *UserUseCaseImpl) GetUserInfo(claims *security.AccessTokenModel) (*GetUserInfoRequest, error) {
	if claims.Role != user.RoleUser {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE003002001",
			"Failed to check role.",
			"Failed to check role.",
		)
	}

	existUser, err := u.userRepo.FindUserById(claims.ID)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE003002002",
			"User not found.",
			err.Error(),
		)
	}

	if existUser.State == user.UserStateSuspended {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE003002003",
			"User suspended.",
			"User suspended.",
		)
	}

	return &GetUserInfoRequest{
		Name:      existUser.Name,
		GroupCode: existUser.GroupCode,
	}, nil
}

// RemoveUserAndInfo implements UserUseCaseInterface.
func (u *UserUseCaseImpl) RemoveUserAndInfo(claims *security.AccessTokenModel) error {
	if claims.Role != user.RoleUser {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE003003001",
			"Failed to check role.",
			"Failed to check role.",
		)
	}

	err := u.userRepo.RemoveUserById(claims.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE003003002",
			"Failed to remove user.",
			"Failed to remove user.",
		)
	}

	u.cardRecordRepo.RemoveDataByUserId(claims.ID)
	u.groupRecordRepo.RemoveDataByUserId(claims.ID)
	u.storyRecordRepo.RemoveDataByUserId(claims.ID)
	u.authRepo.RemoveTokenByUserId(claims.ID)

	return nil
}
