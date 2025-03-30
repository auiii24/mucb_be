package user

import (
	"regexp"
)

func ValidateThaiPhoneNumber(phoneNumber string) (bool, error) {
	phoneRegex := `^\+66[689]\d{8}$`

	matched, err := regexp.MatchString(phoneRegex, phoneNumber)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}
	return true, nil
}
