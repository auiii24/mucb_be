package user

type UserRepository interface {
	CreateUser(user *User) error
	FindUserByPhoneNumber(phoneNumber string) (*User, error)
	FindUserById(id string) (*User, error)
	UpdateUserInfo(id, name, group string) error
}
