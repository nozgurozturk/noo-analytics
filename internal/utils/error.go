package errors

import "net/http"

type ApplicationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

func BadRequest(message string) *ApplicationError {
	return &ApplicationError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func InternalServer(message string) *ApplicationError {
	return &ApplicationError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "server_error",
	}
}

func NotFound(message string) *ApplicationError {
	return &ApplicationError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "server_error",
	}
}

func AlreadyExist(message string) *ApplicationError {
	return &ApplicationError{
		Message: message,
		Status:  http.StatusConflict,
		Error:   "server_error",
	}
}

func Unauthorized(message string) *ApplicationError {
	return &ApplicationError{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

func Forbidden(message string) *ApplicationError {
	return &ApplicationError{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   "forbidden",
	}
}
