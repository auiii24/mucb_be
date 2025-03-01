package auth

type AuthUseCaseInterface interface {
	SignInAdmin(req *SignInAdminRequest, userAgent string) (*SignInAdminOutput, error)
	RenewAdmin(req *RenewAdminRequest, userAgent string) (*RenewAdminOutput, error)
	SignInUser(req *SignInUserRequest) (*SignInUserOutput, error)
	VerifyOtp(req *VerifyOtpRequest) (*VerifyOtpOutput, error)
	RenewUser(req *RenewUserRequest, userAgent string) (*RenewUserOutput, error)
	SignOut(req *SignOutRequest) error
}
