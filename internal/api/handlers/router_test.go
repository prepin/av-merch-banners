package handlers

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestConfigureFieldValidator(t *testing.T) {
	t.Run("json and form tags", func(t *testing.T) {
		handlers := &Handlers{}
		gin.SetMode(gin.TestMode)

		type TestStruct struct {
			JSONField  string `json:"json_field" binding:"required"`
			FormField  string `form:"form_field" binding:"required"`
			EmptyField string
		}

		handlers.configureFieldValidator()

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			test := TestStruct{}
			err := v.Struct(test)

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				for _, e := range validationErrors {
					switch e.Field() {
					case "json_field", "form_field":
						assert.Equal(t, "required", e.Tag())
					default:
						t.Errorf("unexpected validation error for field: %s", e.Field())
					}
				}
			}
		}
	})

	t.Run("invalid validator engine", func(t *testing.T) {
		handlers := &Handlers{}

		original := binding.Validator
		mockVal := &mockValidator{}
		binding.Validator = mockVal

		mockVal.On("Engine").Return(nil)

		defer func() {
			binding.Validator = original
			mockVal.AssertExpectations(t)
		}()

		handlers.configureFieldValidator()
	})
}
