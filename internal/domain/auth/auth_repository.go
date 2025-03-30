package auth

type AuthRepository interface {
	CreateToken(token Token) error
	FindTokenById(id string) (*Token, error)
	RemoveTokenById(id string) error
	RemoveTokenByUserId(id string) error
	FindAllTokenByUserId(id string) (*[]Token, error)
	UpdateTimestampByTokenId(id string) error
}
