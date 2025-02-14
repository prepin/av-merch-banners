package handlers

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"errors"`
}

var NotFoundResponse = ErrorResponse{Error: "not found"}
var InvalidRequestResponse = ErrorResponse{Error: "invalid request"}
var ServerErrorResponse = ErrorResponse{Error: "server error"}
var UnauthorizedResponse = ErrorResponse{Error: "unauthorized"}

func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				return field + " is required"
			default:
				return field + " is invalid"
			}
		}
	}
	return "Invalid input parameters"
}
