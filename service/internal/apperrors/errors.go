package apperrors

import "errors"

// Базовые ошибки приложения
var (
	ErrNewsNotFound = errors.New("news not found")
	ErrInvalidID    = errors.New("invalid id format")
	ErrInvalidBody  = errors.New("invalid request body")
	ErrValidation   = errors.New("validation failed")
)

// AppError - кастомная ошибка с HTTP статусом
type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Конструкторы
func NewBadRequest(message string) *AppError {
	return &AppError{
		Err:        ErrInvalidBody,
		Message:    message,
		StatusCode: 400,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Err:        ErrNewsNotFound,
		Message:    message,
		StatusCode: 404,
	}
}

func NewValidation(message string) *AppError {
	return &AppError{
		Err:        ErrValidation,
		Message:    message,
		StatusCode: 400,
	}
}

func NewInternal(message string) *AppError {
	return &AppError{
		Err:        errors.New("internal error"),
		Message:    message,
		StatusCode: 500,
	}
}
