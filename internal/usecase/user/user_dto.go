package user

type UpdateUserInfoRequest struct {
	Name      string `json:"name" binding:"required,max=128"`
	GroupCode string `json:"group" binding:"max=32"`
}

type UpdateUserInfoOutput struct {
	AccessToken string `json:"accessToken"`
}

type RenewUserRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,max=1024"`
}

type RenewUserOutput struct {
	AccessToken string `json:"accessToken"`
}
