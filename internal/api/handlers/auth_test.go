package handlers

import (
	"av-merch-shop/internal/errs"
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

func TestPostAuth(t *testing.T) {
	mockLogger := slog.Default()
	mockUseCase := &mockAuthUseCase{}
	handler := NewAuthHandler(mockLogger, mockUseCase)

	t.Run("error when signing in", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.NewReader(`{"username": "test", "password": "test"}`)
		c.Request = httptest.NewRequest("POST", "/auth", body)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUseCase.On("SignIn", mock.Anything, "test", "test").
			Return("", errors.New("internal error"))

		handler.PostAuth(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server error")
	})

	t.Run("no access error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := strings.NewReader(`{"username": "test", "password": "wrong"}`)
		c.Request = httptest.NewRequest("POST", "/auth", body)
		c.Request.Header.Set("Content-Type", "application/json")

		mockUseCase.On("SignIn", mock.Anything, "test", "wrong").
			Return("", errs.ErrNoAccess{})

		handler.PostAuth(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "wrong password")
	})
}
