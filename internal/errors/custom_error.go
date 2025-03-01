package errors

import "fmt"

type CustomError struct {
	StatusCode    int    `json:"-"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	SystemMessage string `json:"-"`
}

func (e CustomError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.SystemMessage)
}

func NewCustomError(statusCode int, code string, message, systemMessage string) *CustomError {
	return &CustomError{
		Code:          code,
		Message:       message,
		SystemMessage: systemMessage,
		StatusCode:    statusCode,
	}
}
