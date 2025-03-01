package auth

type OtpAttemptRepository interface {
	CheckOtpRateLimit(phoneNumber string) (bool, error)
	IncrementOtpAttempt(phoneNumber string) error
	ClearOtpAttempts(phoneNumber string) error
}
