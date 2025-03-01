package user

import "mucb_be/internal/infrastructure/security"

type UserUseCaseInterface interface {
	UpdateUserInfo(req *UpdateUserInfoRequest, claims *security.AccessTokenModel) (*UpdateUserInfoOutput, error)
}
