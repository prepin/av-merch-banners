package handlers

import (
	"av-merch-shop/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	logger  *slog.Logger
	usecase *usecase.InfoUseCase
}

func NewInfoHandler(l *slog.Logger, u *usecase.InfoUseCase) *InfoHandler {
	return &InfoHandler{
		logger:  l,
		usecase: u,
	}
}

// возвращает инфу про текущего пользователя
func (h *InfoHandler) GetInfo(c *gin.Context) {
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

	info, err := h.usecase.GetInfo(c.Request.Context(), userIDInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	c.JSON(http.StatusOK, info)

}
