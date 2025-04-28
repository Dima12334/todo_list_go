package errors

import "errors"

var (
	ErrUserNotFound          = errors.New("user doesn't exists")
	ErrTaskNotFound          = errors.New("task not found")
	ErrUserAlreadyExists     = errors.New("user with such email already exists")
	ErrTaskAlreadyExists     = errors.New("task with such title already exists")
	ErrCategoryAlreadyExists = errors.New("category with such title already exists")
)
