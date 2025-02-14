package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostCredit(t *testing.T) {
	mockLogger := slog.Default()
	mockUseCase := &mockCreditUseCase{}
	handler := NewCreditHandler(mockLogger, mockUseCase)

	t.Run("invalid json binding", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.NewReader(`{"invalid": json}`)
		c.Request = httptest.NewRequest("POST", "/credit", body)
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostCredit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("server error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.NewReader(`{"username": "test", "amount": 100}`)
		c.Request = httptest.NewRequest("POST", "/credit", body)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUseCase.On("Credit", mock.Anything, mock.AnythingOfType("*entities.CreditData")).
			Return(nil, errors.New("internal error"))

		handler.PostCredit(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
