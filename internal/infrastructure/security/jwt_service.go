package security

import (
	"errors"
	"fmt"
	"mucb_be/internal/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJWTUnexpectedMethod = errors.New("unexpected signing method")
	ErrJWTMissingIDClaim   = errors.New("invalid token: id claim is missing or not a string")
	ErrJWTMissingRoleClaim = errors.New("invalid token: role claim is missing")
	ErrJWTTokenExpired     = errors.New("token has expired")
)

type JwtServiceInterface interface {
	GenerateAccessToken(id string, role string) (string, error)
	ValidateAccessToken(accessToken string) (AccessTokenModel, error)
}

type AccessTokenModel struct {
	ID   string
	Role string
}

type JwtService struct {
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) JwtServiceInterface {
	return &JwtService{
		cfg: cfg,
	}
}

func (s JwtService) GenerateAccessToken(id string, role string) (string, error) {
	minutes, err := strconv.Atoi(s.cfg.AccessTokenExpiredMinute)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Duration(minutes) * time.Minute).Unix()

	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.AccessTokenKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s JwtService) ValidateAccessToken(accessToken string) (AccessTokenModel, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrJWTUnexpectedMethod, token.Header["alg"])
		}

		return []byte(s.cfg.AccessTokenKey), nil
	})
	if err != nil {
		return AccessTokenModel{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id, ok := claims["id"].(string)
		if !ok {
			return AccessTokenModel{}, ErrJWTMissingIDClaim
		}

		role, ok := claims["role"].(string)
		if !ok {
			return AccessTokenModel{}, ErrJWTMissingRoleClaim
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return AccessTokenModel{}, ErrJWTTokenExpired
		}

		// Convert exp to Unix timestamp and compare with current time
		if int64(exp) < time.Now().Unix() {
			return AccessTokenModel{}, ErrJWTTokenExpired
		}

		return AccessTokenModel{
			ID:   id,
			Role: role,
		}, nil
	} else {
		return AccessTokenModel{}, err
	}
}
