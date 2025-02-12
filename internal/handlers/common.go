package handlers

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"errors"`
}

var NotFoundResponse = ErrorResponse{Error: "not found"}
var InvalidRequestResponse = ErrorResponse{Error: "invalid request"}
var ServerErrorResponse = ErrorResponse{Error: "server error"}

func formatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
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
