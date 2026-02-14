package domain

import (
	"errors"
	"fmt"
)

var (
	ErrConflict       = errors.New("conflict: item already exists")
	ErrNotFound       = errors.New("not found: resource does not exist")
	ErrUnauthorized   = errors.New("unauthorized: authentication required")
	ErrForbidden      = errors.New("forbidden: insufficient permissions")
	ErrInvalidInput   = errors.New("invalid input: check provided data")
	ErrInternal       = errors.New("internal error: something went wrong")
	ErrAlreadyExists  = errors.New("already exists")
	ErrNotImplemented = errors.New("not implemented")
	ErrTimeout        = errors.New("operation timed out")
)

type ErrorCode string

const (
	CodeConflict       ErrorCode = "CONFLICT"
	CodeNotFound       ErrorCode = "NOT_FOUND"
	CodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	CodeForbidden      ErrorCode = "FORBIDDEN"
	CodeInvalidInput   ErrorCode = "INVALID_INPUT"
	CodeInternal       ErrorCode = "INTERNAL"
	CodeAlreadyExists  ErrorCode = "ALREADY_EXISTS"
	CodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
	CodeTimeout        ErrorCode = "TIMEOUT"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return errors.Is(e.Err, target)
}

func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func ErrConflictWith(cause error, msg string) *AppError {
	return &AppError{Code: CodeConflict, Message: msg, Err: cause}
}

func ErrNotFoundWith(cause error, msg string) *AppError {
	return &AppError{Code: CodeNotFound, Message: msg, Err: cause}
}

func ErrUnauthorizedWith(cause error, msg string) *AppError {
	return &AppError{Code: CodeUnauthorized, Message: msg, Err: cause}
}

func ErrForbiddenWith(cause error, msg string) *AppError {
	return &AppError{Code: CodeForbidden, Message: msg, Err: cause}
}

func ErrInvalidInputWith(cause error, msg string) *AppError {
	return &AppError{Code: CodeInvalidInput, Message: msg, Err: cause}
}

func ErrInternalWith(cause error, msg string) *AppError {
	return &AppError{Code: CodeInternal, Message: msg, Err: cause}
}

func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict) || isCode(err, CodeConflict)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || isCode(err, CodeNotFound)
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized) || isCode(err, CodeUnauthorized)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden) || isCode(err, CodeForbidden)
}

func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput) || isCode(err, CodeInvalidInput)
}

func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal) || isCode(err, CodeInternal)
}

func isCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
