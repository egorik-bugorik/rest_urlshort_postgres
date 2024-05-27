package resp

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`

	Error string `json:"error,omitempty"`
}

const (
	ResponseOk    = "OK"
	ResponseError = "Error"
)

func Error(s string) *Response {
	return &Response{
		Status: ResponseError,
		Error:  s,
	}
}
func OK(s string) *Response {
	return &Response{
		Status: ResponseOk,
	}
}

func ValidateError(err validator.ValidationErrors) *Response {

	var errMsg []string

	for _, fieldError := range err {
		switch fieldError.Tag() {
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("Url field is wrong ::: %s", fieldError.Field()))
		case "required":
			errMsg = append(errMsg, fmt.Sprintf(" field is required ::: %s", fieldError.Field()))
		default:

			errMsg = append(errMsg, fmt.Sprintf("field is wrong ::: %s", fieldError.Field()))

		}

	}

	return Error(strings.Join(errMsg, ", "))

}
