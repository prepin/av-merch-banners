package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInfo(t *testing.T) {
	mockLogger := slog.Default()
	mockUseCase := &mockInfoUseCase{}
	handler := NewInfoHandler(mockLogger, mockUseCase)

	t.Run("userID type assertion failure", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
		c.Set("userID", "not an int")

		handler.GetInfo(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server error")
	})

	t.Run("usecase error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
		c.Set("userID", 1)

		mockUseCase.On("GetInfo", mock.Anything, 1).
			Return(nil, errors.New("internal error"))

		handler.GetInfo(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("missing userID", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/info", http.NoBody)

		handler.GetInfo(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
