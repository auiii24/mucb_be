package auth

type AuthRepository interface {
	CreateToken(token Token) error
	FindTokenById(id string) (*Token, error)
	RemoveTokenById(id string) error
}
