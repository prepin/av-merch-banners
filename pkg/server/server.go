package server

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/handlers"
	"av-merch-shop/pkg/auth"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func New(cfg *config.Config, handlers *handlers.Handlers) *Server {
	router := gin.Default()
	jwtService := auth.NewJWTService(cfg)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost" + cfg.Server.Port[1:],
	}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "User-Agent"}

	router.Use(cors.New(config))
	handlers.RegisterRoutes(router, jwtService)

	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		},
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.cfg.Logger.Error("Failed to launch server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return s.Shutdown()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	select {
	case <-ctx.Done():
		s.cfg.Logger.Warn("timeout of 5 seconds.")
	default:
		s.cfg.Logger.Info("Server shutdown done.")
	}

	return nil
}
