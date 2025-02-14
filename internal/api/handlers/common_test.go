package handlers

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestFormatValidationError(t *testing.T) {
    t.Run("non-validation error", func(t *testing.T) {
        err := errors.New("random error")
        result := formatValidationError(err)
        assert.Equal(t, "Invalid input parameters", result)
    })

    t.Run("unknown validation tag", func(t *testing.T) {
        validate := validator.New()
        type Test struct {
            Field string `validate:"min=5"`
        }
        var test Test
        err := validate.Struct(test)
        result := formatValidationError(err)
        assert.Equal(t, "Field is invalid", result)
    })
}
