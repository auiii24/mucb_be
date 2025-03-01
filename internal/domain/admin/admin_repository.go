package admin

type AdminRepository interface {
	CreateAdmin(admin *Admin) error
	FindAdminByEmail(email string) (*Admin, error)
	FindAdminById(id string) (*Admin, error)

	// CreateUser(user *User) error

	// FindUserByPhone(phone string) (*User, error)
	// StoreToken(token *Token) error
	// FindToken(userID string) (*Token, error)
	// StoreOTP(otp *OTP) error
	// VerifyOTP(phone, code string) (*OTP, error)
}
