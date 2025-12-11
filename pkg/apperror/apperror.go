package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrType string

const (
	NotFoundError       ErrType = "NotFound"
	InternalServerError ErrType = "InternalServerError"
	BadRequestError     ErrType = "BadRequest"
	InvalidDataError    ErrType = "InvalidData"
	UnauthorizedError   ErrType = "Unauthorized"
	ConflictError       ErrType = "Conflict"

	unexpectedErrorMessage = "something went wrong"
)

type AppError struct {
	Err     error   `json:"error"`
	Type    ErrType `json:"type"`
	Message string  `json:"message,omitempty"`
	Code    string  `json:"code"`
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, _ := json.Marshal(e)
	return marshal
}

func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{Err: err, Type: InternalServerError, Message: unexpectedErrorMessage}
}

func NewBadRequestError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: BadRequestError, Message: message, Code: code}
}

func NewInvalidDataError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: InvalidDataError, Message: message, Code: code}
}

func NewUnauthorizedError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: UnauthorizedError, Message: message, Code: code}
}

func NewNotFoundError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: NotFoundError, Message: message, Code: code}
}

func NewConflictError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: ConflictError, Message: message, Code: code}
}

func GetErrorByHttpStatus(status int, message, code string) error {
	switch status {
	case http.StatusBadRequest:
		return NewBadRequestError(message, code)
	case http.StatusUnauthorized:
		return NewUnauthorizedError(message, code)
	case http.StatusUnprocessableEntity:
		return NewInvalidDataError(message, code)
	case http.StatusNotFound:
		return NewNotFoundError(message, code)
	default:
		return NewInternalError(errors.New(message))
	}
}

func GetHttpStatusByErrorType(errType ErrType) int {
	status := http.StatusInternalServerError
	switch errType {
	case InvalidDataError:
		status = http.StatusUnprocessableEntity
	case NotFoundError:
		status = http.StatusNotFound
	case UnauthorizedError:
		status = http.StatusUnauthorized
	case BadRequestError:
		status = http.StatusBadRequest
	}
	return status
}
