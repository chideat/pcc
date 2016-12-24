package error

import (
	"fmt"
)

type UserError struct {
	code    string
	message string
}

func (err *UserError) Error() string {
	return fmt.Sprintf("%s %s", err.Code)
}

func (err *UserError) Code() string {
	return err.code
}

func (err *UserError) Message() string {
	return err.message
}

func NewUserError(code, message string) *UserError {
	return &UserError{code: code, message: message}
}
