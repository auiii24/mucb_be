package user

import (
	"errors"
	"regexp"
)

func ValidateThaiPhoneNumber(phoneNumber string) error {
	phoneRegex := `^(0[689]\d{8}|\+66[689]\d{8})$`

	matched, err := regexp.MatchString(phoneRegex, phoneNumber)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid Thai phone number format")
	}
	return nil
}
