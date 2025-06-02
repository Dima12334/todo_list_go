package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	customErrors "todo_list_go/pkg/errors"
)

func BindAndValidateJSON[T any](c *gin.Context, input *T) (bool, any) {
	var errResponseBody any

	rawData, err := c.GetRawData()
	if err != nil {
		errResponseBody = "Invalid request body"
		return false, errResponseBody
	}

	// Unmarshal to catch JSON syntax/type errors
	decoder := json.NewDecoder(bytes.NewReader(rawData))
	if err = decoder.Decode(input); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			errResponseBody = fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxErr.Offset)
		case errors.As(err, &unmarshalTypeErr):
			errResponseBody = fmt.Sprintf("Field '%s' expects a %s but got %s", unmarshalTypeErr.Field, unmarshalTypeErr.Type, unmarshalTypeErr.Value)
		case errors.Is(err, io.EOF):
			errResponseBody = "Request body cannot be empty"
		default:
			errResponseBody = err.Error()
		}
		return false, errResponseBody
	}

	// Restore body for Gin
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	//Bind & validate using Gin
	if err = c.ShouldBindJSON(input); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, *input)
		if out != nil {
			return false, out
		}

		return false, err.Error()
	}

	return true, errResponseBody
}
