package responses

import (
	"net/http"

	"pvpc-backend/internal/domain/errors"
)

type APIErrorResponse struct {
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func NewAPIErrorResponse(err error) (int, APIErrorResponse) {

	errorCode := getErrorCode(err)
	statusCode := mapErrorToStatusCode(err)

	return statusCode, APIErrorResponse{
		ErrorCode:  errorCode,
		Message:    errors.ErrorWithoutCode(err),
		StatusCode: statusCode,
	}
}

func getErrorCode(err error) string {
	switch errors.Code(err) {
	case "", errors.PersistenceError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return string(errors.Code(err))
	}
}

func mapErrorToStatusCode(err error) int {
	switch errors.Code(err) {
	case errors.InvalidPricesID, errors.InvalidZoneID:
		return http.StatusBadRequest
	case errors.ZoneNotFound:
		return http.StatusNotFound
	case errors.PersistenceError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
