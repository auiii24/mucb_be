package auth

type OtpRepository interface {
	CreateOtp(otp *Otp) error
	FindLatestOtpByPhoneNumber(phoneNumber string) (*Otp, error)
	MarkOtpAsUsedById(id string) error
	IncrementOtpAttemptsById(id string) error
}
