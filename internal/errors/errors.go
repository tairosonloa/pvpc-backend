package errors

import (
	"fmt"
)

type ErrorCode string

const (
	InvalidPricesID     ErrorCode = "INVALID_PRICES_ID"
	InvalidPricesZoneID ErrorCode = "INVALID_PRICES_ZONE_ID"
	PersistenceError    ErrorCode = "PERSISTENCE_ERROR"
	PricesZoneNotFound  ErrorCode = "PRICES_ZONE_NOT_FOUND"
)

type domainError struct {
	error
	errorCode ErrorCode
}

func (e domainError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorCode, e.error.Error())
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

func WrapIntoDomainError(err error, errorCode ErrorCode, msg string) error {
	return domainError{
		error:     fmt.Errorf("%s: [%w]", msg, err),
		errorCode: errorCode,
	}
}

func NewDomainError(errorCode ErrorCode, format string, args ...interface{}) error {
	return domainError{
		error:     fmt.Errorf(format, args...),
		errorCode: errorCode,
	}
}
