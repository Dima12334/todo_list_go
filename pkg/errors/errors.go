package errors

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"reflect"
	"strings"
)

var (
	ErrUserNotFound          = errors.New("user doesn't exists")
	ErrTaskNotFound          = errors.New("task not found")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrUserAlreadyExists     = errors.New("user with such email already exists")
	ErrTaskAlreadyExists     = errors.New("task with such title already exists")
	ErrCategoryAlreadyExists = errors.New("category with such title already exists")
	ErrNoUpdateFields        = errors.New("no fields specified for update")
)

func IsDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // Unique violation code
	}
	return false
}

func FormatValidationErrorOutput(err error, inputObj any) map[string]string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string)
		t := reflect.TypeOf(inputObj)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		for _, fe := range ve {
			if field, ok := t.FieldByName(fe.StructField()); ok {
				jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
				if jsonTag == "-" || jsonTag == "" {
					jsonTag = strings.ToLower(fe.Field())
				}
				out[jsonTag] = ValidationErrorToText(fe)
			}
		}
		return out
	}
	return nil
}

func ValidationErrorToText(fe validator.FieldError) string {
	fieldKind := fe.Kind()

	switch fe.Tag() {
	case "uuid":
		return "must be a valid UUID"
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		if fieldKind == reflect.String {
			return fmt.Sprintf("must be at least %s characters", fe.Param())
		}
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		if fieldKind == reflect.String {
			return fmt.Sprintf("must be at most %s characters", fe.Param())
		}
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	default:
		return "is not valid"
	}
}
