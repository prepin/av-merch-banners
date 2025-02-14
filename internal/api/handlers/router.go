package handlers

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/api/middleware"
	"av-merch-shop/internal/usecase"
	"av-merch-shop/pkg/auth"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	timeout "github.com/vearne/gin-timeout"
)

type Handlers struct {
	config *config.Config
	Auth   *AuthHandler
	Credit *CreditHandler
	Send   *SendCoinHandler
	Order  *OrderHandler
	Info   *InfoHandler
}

func NewHandlers(cfg *config.Config, usecases usecase.Usecases) *Handlers {
	return &Handlers{
		config: cfg,
		Auth:   NewAuthHandler(cfg.Logger, usecases.AuthUseCase),
		Credit: NewCreditHandler(cfg.Logger, usecases.CreditUseCase),
		Send:   NewSendCoinHandler(cfg.Logger, usecases.SendCoinUseCase),
		Order:  NewOrderHandler(cfg.Logger, usecases.OrderUseCase),
		Info:   NewInfoHandler(cfg.Logger, usecases.InfoUseCase),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine, jwtService *auth.JWTService) {

	h.configureFieldValidator()

	api := r.Group("/api")
	api_v1 := r.Group("/api/v1")

	registerRoutes := func(groups ...*gin.RouterGroup) {
		for _, g := range groups {

			// авторизация (и автосоздание юзера)
			g.POST("/auth", h.Auth.PostAuth)

			timed := g.Group("")

			timeoutMsg := `{"error":"timeout"}`
			timed.Use(timeout.Timeout(
				timeout.WithTimeout(time.Duration(h.config.Server.RequestTimeout)*time.Millisecond),
				timeout.WithErrorHttpCode(http.StatusRequestTimeout),
				timeout.WithDefaultMsg(timeoutMsg),
				timeout.WithGinCtxCallBack(func(c *gin.Context) {
					h.config.Logger.Warn("Timeout", "url", c.Request.URL.String())
				})))

			// healthcheck эндпойнты
			{
				timed.GET("/ping", GetPingHandler)
				timed.GET("/teapot", GetTeapotHandler)
				timed.GET("/sleep", GetSleepHandler)
			}

			protected := timed.Group("")
			protected.Use(middleware.AuthMiddleware(jwtService))
			{
				// отправка монет другому пользователю
				protected.POST("/sendCoin", h.Send.PostSendCoin)

				// покупка товара (второй роут для отлова запросов без item вообще)
				protected.POST("/buy/:item", h.Order.PostOrder)
				protected.POST("/buy/", h.Order.PostOrder)

				// тоже покупка, но по GET, для совместимости
				protected.GET("/buy/:item", h.Order.PostOrder)
				protected.GET("/buy/", h.Order.PostOrder)

				// Информация о текущем пользователе
				protected.GET("/info", h.Info.GetInfo)
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
