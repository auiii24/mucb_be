package auth

import "mucb_be/internal/domain/user"

type RefreshTokenPayload struct {
	User  string `json:"user"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

type SignInAdminRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignInAdminOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RenewAdminRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,max=1024"`
}

type RenewAdminOutput struct {
	AccessToken string `json:"accessToken"`
}

type SignInUserRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,max=10"`
}

type SignInUserOutput struct {
	Token string `json:"token"`
}

type VerifyOtpRequest struct {
	Token string `json:"token" binding:"required"`
	Code  string `json:"code" binding:"required,max=6"`
}

type VerifyOtpOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Status       string `json:"status"`
}

type RenewUserRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,max=1024"`
}

type RenewUserOutput struct {
	AccessToken string `json:"accessToken"`
}

type SignOutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,max=1024"`
}

func (req *SignInUserRequest) Validate() error {
	if err := user.ValidateThaiPhoneNumber(req.PhoneNumber); err != nil {
		return err
	}
	return nil
}
