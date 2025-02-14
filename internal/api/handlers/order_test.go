package handlers

import (
	"av-merch-shop/internal/entities"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostOrder(t *testing.T) {
	mockLogger := slog.Default()
	mockUseCase := &mockOrderUseCase{}
	handler := NewOrderHandler(mockLogger, mockUseCase)

	t.Run("server error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/buy/item1", http.NoBody)
		c.Set("userID", 1)
		c.AddParam("item", "item1")

		mockUseCase.On("Buy", mock.Anything, &entities.OrderRequest{
			UserID:   1,
			ItemName: "item1",
		}).Return(errors.New("internal server error"))

		handler.PostOrder(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "server error")
		mockUseCase.AssertExpectations(t)
	})

	t.Run("wrong method", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/buy/item1", http.NoBody)
		c.Set("userID", 1)
		c.AddParam("item", "item1")

		mockUseCase.On("Buy", mock.Anything, &entities.OrderRequest{
			UserID:   1,
			ItemName: "item1",
		}).Return(nil)

		handler.PostOrder(c)

		assert.Equal(t, "true", w.Header().Get("Deprecation"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("empty item param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/buy/", http.NoBody)
		c.AddParam("item", "/")

		handler.PostOrder(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "item is required")
	})

	t.Run("missing userID", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/buy/item1", http.NoBody)
		c.AddParam("item", "item1")

		handler.PostOrder(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("userID type assertion failure", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/buy/item1", http.NoBody)
		c.Set("userID", "not an int")
		c.AddParam("item", "item1")

		handler.PostOrder(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}
