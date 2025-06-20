package errors

import "net/http"

type ErrorType uint32

const (
	BadRequest       = ErrorType(http.StatusBadRequest)
	NotFound         = ErrorType(http.StatusNotFound)
	Teapot           = ErrorType(http.StatusTeapot)
	InternalServer   = ErrorType(http.StatusInternalServerError)
	Forbidden        = ErrorType(http.StatusForbidden)
	Unauthorized     = ErrorType(http.StatusUnauthorized)
	ClientReject     = ErrorType(http.StatusNotAcceptable)
	BadGateway       = ErrorType(http.StatusBadGateway)
	ContextCancelled = ErrorType(http.StatusRequestTimeout)
)
