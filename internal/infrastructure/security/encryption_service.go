package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mucb_be/internal/config"
	"strings"
	"unicode"
)

type EncryptionServiceInterface interface {
	EncryptRefreshToken(plaintext string) (string, error)
	DecryptRefreshToken(ciphertext string) (string, error)
	EncryptData(plaintext string) (string, error)
	DecryptData(ciphertext string) (string, error)
}

type EncryptionService struct {
	cfg *config.Config
}

func NewEncryptionService(cfg *config.Config) EncryptionServiceInterface {
	return &EncryptionService{
		cfg: cfg,
	}
}

func (s *EncryptionService) EncryptRefreshToken(plaintext string) (string, error) {
	return s.encryptDataWithKey(plaintext, s.cfg.RefreshTokenKey)
}

func (s *EncryptionService) DecryptRefreshToken(ciphertext string) (string, error) {
	return s.decryptDataWithKey(ciphertext, s.cfg.RefreshTokenKey)
}

func (s *EncryptionService) EncryptData(plaintext string) (string, error) {
	return s.encryptDataWithKey(plaintext, s.cfg.EncryptionKey)
}
func (s *EncryptionService) DecryptData(ciphertext string) (string, error) {
	return s.decryptDataWithKey(ciphertext, s.cfg.EncryptionKey)
}

func (s *EncryptionService) encryptDataWithKey(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Generate a random initialization vector (IV).
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	// Encode the ciphertext as base64 for easy storage and transmission.
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *EncryptionService) decryptDataWithKey(ciphertext, key string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// The IV is the first block of the ciphertext.
	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	decryptedText := string(ciphertextBytes)
	decryptedText = strings.TrimSpace(decryptedText)
	decryptedText = strings.TrimFunc(decryptedText, func(r rune) bool {
		return !unicode.IsPrint(r) // Remove non-printable characters
	})

	return decryptedText, nil
}
