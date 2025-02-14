package handlers

import (
	"av-merch-shop/internal/entities"
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

func TestPostSendCoin(t *testing.T) {
	mockLogger := slog.Default()
	mockUseCase := &mockSendCoinUseCase{}
	handler := NewSendCoinHandler(mockLogger, mockUseCase)

	t.Run("internal server error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin",
			strings.NewReader(`{"toUser": "test", "amount": 100}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", 1)

		mockUseCase.On("Send", mock.Anything, &entities.TransferData{
			SenderID:  1,
			Recipient: "test",
			Amount:    100,
		}).Return(errors.New("internal server error"))

		handler.PostSendCoin(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server error")
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(`{invalid json}`))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostSendCoin(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing userID", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin",
			strings.NewReader(`{"toUser": "test", "amount": 100}`))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.PostSendCoin(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("userID type assertion failure", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin",
			strings.NewReader(`{"toUser": "test", "amount": 100}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "not an int")

		handler.PostSendCoin(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin",
			strings.NewReader(`{"toUser": "nonexistent", "amount": 100}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", 1)

		mockUseCase.On("Send", mock.Anything, &entities.TransferData{
			SenderID:  1,
			Recipient: "nonexistent",
			Amount:    100,
		}).Return(errs.ErrNotFoundError)

		handler.PostSendCoin(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
