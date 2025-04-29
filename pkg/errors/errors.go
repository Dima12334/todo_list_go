package errors

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

var (
	ErrUserNotFound          = errors.New("user doesn't exists")
	ErrTaskNotFound          = errors.New("task not found")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrUserAlreadyExists     = errors.New("user with such email already exists")
	ErrTaskAlreadyExists     = errors.New("task with such title already exists")
	ErrCategoryAlreadyExists = errors.New("category with such title already exists")
)

func IsDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // Unique violation code
	}
	return false
}

func ValidationErrorToText(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	default:
		return "is not valid"
	}
}
