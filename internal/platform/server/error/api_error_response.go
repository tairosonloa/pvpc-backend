package error

import (
	pvpc "go-pvpc/internal"
	"net/http"
)

type APIErrorResponse struct {
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func NewAPIErrorResponse(err error) APIErrorResponse {

	errorCode := DeduceErrorCode(err)
	statusCode := DeduceStatusCode(err)

	return APIErrorResponse{
		ErrorCode:  errorCode,
		Message:    err.Error(),
		StatusCode: statusCode,
	}
}

func DeduceStatusCode(err error) int {
	switch err {
	case pvpc.ErrInvalidPricesID, pvpc.ErrInvalidPricesZoneID:
		return http.StatusBadRequest
	case pvpc.ErrPricesZoneNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func DeduceErrorCode(err error) string {
	switch err {
	case pvpc.ErrInvalidPricesID:
		return "INVALID_PRICES_ID"
	case pvpc.ErrInvalidPricesZoneID:
		return "INVALID_PRICES_ZONE_ID"
	case pvpc.ErrPricesZoneNotFound:
		return "PRICES_ZONE_NOT_FOUND"
	default:
		return "INTERNAL_SERVER_ERROR"
	}
}
