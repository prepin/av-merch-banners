package handlers

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"av-merch-shop/internal/usecase"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreditHandler struct {
	logger  *slog.Logger
	usecase *usecase.CreditUseCase
}

type CreditParams struct {
	Username string `json:"username" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
}

type CreditResponse struct {
	NewAmount   int       `json:"new_amount"`
	ReferenceId uuid.UUID `json:"reference_id"`
}

func NewCreditHandler(l *slog.Logger, u *usecase.CreditUseCase) *CreditHandler {
	return &CreditHandler{
		logger:  l,
		usecase: u,
	}
}

// Добавляет нужное количество монет указанному пользователю. Если количество монет
// отрицательное, то монеты спишутся (но не ниже нуля).
func (h *CreditHandler) PostCredit(c *gin.Context) {
	var params *CreditParams

	if err := c.ShouldBindJSON(&params); err != nil {
		h.logger.Debug("Failed parsing credit request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: formatValidationError(err)})
		return
	}

	tr, err := h.usecase.Credit(c.Request.Context(), &entities.CreditData{
		Username: params.Username,
		Amount:   params.Amount,
	})

	if err != nil {
		if errors.Is(err, errs.ErrNotFound{}) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not exists"})
			return
		}
		h.logger.Error("Crediting user failed", "error", err)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
	}

	c.JSON(http.StatusCreated, CreditResponse{
		NewAmount:   tr.NewAmount,
		ReferenceId: tr.ReferenceID,
	})
}
