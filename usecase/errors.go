package usecase

import (
	"errors"
	"fmt"
)

type ErrorType uint

const (
	NoType ErrorType = iota
	BadRequest
	InternalError
	NotFound
)

type CustomError struct {
	errorType     ErrorType
	originalError error
}

type errorContext struct {
	Field   string
	Message string
}

func (errorType ErrorType) New(msg string) error {
	return CustomError{errorType: errorType, originalError: errors.New(msg)}
}

func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return CustomError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

func (errorType ErrorType) Wrap(err error) error {
	return CustomError{errorType: errorType, originalError: err}
}

func (error CustomError) Error() string {
	return error.originalError.Error()
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(CustomError); ok {
		return customErr.errorType
	}
	return NoType
}
