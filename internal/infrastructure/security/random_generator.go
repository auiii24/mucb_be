package security

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRefCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	refCode := make([]byte, 4)
	for i := range refCode {
		refCode[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(refCode)
}

func GenerateOtpCode() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := seededRand.Intn(1000000)
	return fmt.Sprintf("%06d", otp)
}
