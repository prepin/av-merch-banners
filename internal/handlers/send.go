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

type SendCoinsHandler struct {
	logger  *slog.Logger
	usecase *usecase.SendCoinsUseCase
}

type SendCoinParams struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

func NewSendCoinHandler(l *slog.Logger, u *usecase.SendCoinsUseCase) *SendCoinsHandler {
	return &SendCoinsHandler{
		logger:  l,
		usecase: u,
	}
}

// Отправляет монеты другому пользователю
func (h *SendCoinsHandler) PostSendCoins(c *gin.Context) {
	var params *SendCoinParams

	if err := c.ShouldBindJSON(&params); err != nil {
		h.logger.Debug("Failed parsing transfer request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: formatValidationError(err)})
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

	err := h.usecase.Send(c.Request.Context(), &entities.TransferData{
		SenderID:  userIDInt,
		Recipient: params.ToUser,
		Amount:    params.Amount,
	})

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInsufficientFundsError):
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: "insufficient funds"})
		case errors.Is(err, errs.ErrIncorrectAmountError):
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: "incorrect amount"})
		case errors.Is(err, errs.ErrNotFoundError):
			c.JSON(http.StatusNotFound, NotFoundResponse)
		default:
			c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		}
		return
	}
	c.Status(http.StatusOK)
}
