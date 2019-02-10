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

type customError struct {
	errorType     ErrorType
	originalError error
}

type errorContext struct {
	Field   string
	Message string
}

func (errorType ErrorType) New(msg string) error {
	return customError{errorType: errorType, originalError: errors.New(msg)}
}

func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(msg, args...)}
}

func (errorType ErrorType) FromError(err error) error {
	return customError{errorType: errorType, originalError: err}
}

func (error customError) Error() string {
	return error.originalError.Error()
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}
	return NoType
}
