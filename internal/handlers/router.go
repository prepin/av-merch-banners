package handlers

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/infrastructure/middleware"
	"av-merch-shop/internal/usecase"
	"av-merch-shop/pkg/auth"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	config *config.Config
	Auth   *AuthHandler
	Credit *CreditHandler
	Send   *SendCoinHandler
}

func NewHandlers(cfg *config.Config, usecases usecase.Usecases) *Handlers {
	return &Handlers{
		config: cfg,
		Auth:   NewAuthHandler(cfg.Logger, usecases.AuthUseCase),
		Credit: NewCreditHandler(cfg.Logger, usecases.CreditUseCase),
		Send:   NewSendCoinHandler(cfg.Logger, usecases.SendCoinUseCase),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine, jwtService *auth.JWTService) {

	h.configureFieldValidator()

	api := r.Group("/api")
	api_v1 := r.Group("/api/v1")

	registerRoutes := func(groups ...*gin.RouterGroup) {
		for _, g := range groups {
			// g.Use(
			// 	middleware.TimeoutMiddleware(
			// 		time.Duration(h.config.Server.RequestTimeout) * time.Millisecond,
			// 	),
			// )

			// healthcheck эндпойнты
			{
				g.GET("/ping", GetPingHandler)
				g.GET("/teapot", GetTeapotHandler)
				g.GET("/sleep", GetSleepHandler)
			}

			// авторизация (и автосоздание юзера)
			g.POST("/auth", h.Auth.PostAuth)

			protected := g.Group("")
			protected.Use(middleware.AuthMiddleware(jwtService))
			{
				// отправка монет другому пользователю
				protected.POST("/sendCoin", h.Send.PostSendCoin)
			}

			{ // роуты для пользователей с ролью админа,
				admin := protected.Group("")
				admin.Use(middleware.AdminOnly())
				{
					// начислить деньги пользователю
					admin.POST("/credit", h.Credit.PostCredit)
				}
			}

		}
	}

	registerRoutes(api, api_v1)
}

// Конфигурим валидатор Джина для того, чтобы он брал имя поля из тегов структуры, а не
// её полей. Это нужно для форматирования ошибок.
func (*Handlers) configureFieldValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "-" && name != "" {
				return name
			}
			if name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]; name != "-" && name != "" {
				return name
			}
			return fld.Name
		})
	}
}
