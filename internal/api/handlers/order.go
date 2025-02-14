package handlers

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"av-merch-shop/internal/usecase"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	logger  *slog.Logger
	usecase usecase.OrderUseCase
}

type OrderParams struct {
	Item string `json:"item" binding:"required"`
}

func NewOrderHandler(l *slog.Logger, u usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		logger:  l,
		usecase: u,
	}
}

// Оформляет заказ на вещь для текущего пользователя
func (h *OrderHandler) PostOrder(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.Header("Deprecation", "true")
	}

	itemParam := c.Param("item")

	if itemParam == "" || itemParam == "/" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "item is required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		h.logger.Error("Failed to convert userID to int")
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	err := h.usecase.Buy(c.Request.Context(), &entities.OrderRequest{
		UserID:   userIDInt,
		ItemName: itemParam,
	})

	if err != nil {
		if errors.Is(err, errs.ErrNotFoundError) {
			c.JSON(http.StatusNotFound, NotFoundResponse)
			return
		}

		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	c.Status(http.StatusOK)
}
