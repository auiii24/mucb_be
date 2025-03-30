package user

import "mucb_be/internal/infrastructure/security"

type UserUseCaseInterface interface {
	UpdateUserInfo(req *UpdateUserInfoRequest, claims *security.AccessTokenModel) (*UpdateUserInfoOutput, error)
	GetUserInfo(claims *security.AccessTokenModel) (*GetUserInfoRequest, error)
	RemoveUserAndInfo(claims *security.AccessTokenModel) error
}
