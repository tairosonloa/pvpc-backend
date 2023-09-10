package errors

import (
	stdErrors "errors"
	"fmt"
)

type ErrorCode string

const (
	InvalidPricesID  ErrorCode = "INVALID_PRICES_ID"
	InvalidTime      ErrorCode = "INVALID_TIME"
	InvalidZoneID    ErrorCode = "INVALID_PRICES_ZONE_ID"
	PersistenceError ErrorCode = "PERSISTENCE_ERROR"
	ZoneNotFound     ErrorCode = "PRICES_ZONE_NOT_FOUND"
)

type domainError struct {
	error
	errorCode ErrorCode
}

func (e domainError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorCode, e.error.Error())
}

func ErrorWithoutCode(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(domainError); ok {
		return e.error.Error()
	}

	return err.Error()
}

func Unwrap(err error) error {
	if e, ok := err.(domainError); ok {
		return stdErrors.Unwrap(e.error)
	}

	return stdErrors.Unwrap(err)
}

func Code(err error) ErrorCode {
	if err == nil {
		return ""
	}

	if e, ok := err.(domainError); ok {
		return e.errorCode
	}

	return ""
}

func NewDomainError(errorCode ErrorCode, format string, args ...interface{}) error {
	return domainError{
		error:     fmt.Errorf(format, args...),
		errorCode: errorCode,
	}
}

func WrapIntoDomainError(err error, errorCode ErrorCode, msg string) error {
	return domainError{
		error:     fmt.Errorf("%s: [%w]", msg, err),
		errorCode: errorCode,
	}
}
