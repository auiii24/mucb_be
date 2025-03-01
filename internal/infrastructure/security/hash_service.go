package security

import "golang.org/x/crypto/bcrypt"

type HashServiceInterface interface {
	HashPassword(password string) (string, error)
	CheckHashPassword(cPassword string, hPassword string) bool
}

type HashService struct{}

func NewHashService() HashServiceInterface {
	return &HashService{}
}

func (h *HashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *HashService) CheckHashPassword(cPassword string, hPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hPassword), []byte(cPassword))
	return err == nil
}
