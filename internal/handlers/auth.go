package handlers

import (
	"av-merch-shop/internal/errs"
	"av-merch-shop/internal/usecase"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	logger  *slog.Logger
	usecase *usecase.AuthUseCase
}

type AuthPostParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthTokenResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(l *slog.Logger, uc *usecase.AuthUseCase) *AuthHandler {
	handler := &AuthHandler{
		logger:  l,
		usecase: uc,
	}
	return handler
}

func (h *AuthHandler) PostAuth(c *gin.Context) {
	var params *AuthPostParams
	if err := c.ShouldBindJSON(&params); err != nil {
		h.logger.Debug("Failed parsing auth request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: formatValidationError(err)})
		return
	}

	token, err := h.usecase.SignIn(params.Username, params.Password)
	if err != nil {
		if errors.Is(err, errs.ErrNoAccess{}) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "wrong password"})
			return
		}
		h.logger.Error("Authorization failed", "error", err)
		c.JSON(http.StatusInternalServerError, ServerErrorResponse)
		return
	}

	c.JSON(http.StatusOK, AuthTokenResponse{Token: token})
}
