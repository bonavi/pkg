package errors

import "net/http"

type ErrorType uint8

const (
	BadRequest ErrorType = iota + 1
	NotFound
	Teapot
	InternalServer
	Forbidden
	Unauthorized
	Timeout
	BadGateway
)

var errorTypeToHTTPCode = map[ErrorType]int{
	BadRequest:     http.StatusBadRequest,
	NotFound:       http.StatusNotFound,
	Teapot:         http.StatusTeapot,
	InternalServer: http.StatusInternalServerError,
	Forbidden:      http.StatusForbidden,
	Unauthorized:   http.StatusUnauthorized,
	Timeout:        http.StatusRequestTimeout,
	BadGateway:     http.StatusBadGateway,
}

func (et ErrorType) HTTPCode() int {

	httpCode, ok := errorTypeToHTTPCode[et]
	if !ok {
		return http.StatusInternalServerError
	}

	return httpCode
}
