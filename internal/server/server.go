package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	conf "san/internal/config"
	"san/internal/db"
	"san/internal/handler"
	"san/internal/middleware"
	"san/internal/router"
	"san/pkg/logger"
	"san/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer   *http.Server
	db           *db.Database
	log          logger.Logger
	userHandler  *handler.UserHandler
	postHandler  *handler.PostHandler
	tokenManager token.TokenManager
}

func NewServer(cfg conf.Config, database *db.Database, userHandler *handler.UserHandler, postHandler *handler.PostHandler, tokenManager token.TokenManager, log logger.Logger) *Server {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(log))

	s := &Server{
		db:           database,
		log:          log,
		userHandler:  userHandler,
		postHandler:  postHandler,
		tokenManager: tokenManager,
	}

	r.GET("/health", s.handleHealth)

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.SetupRoutes(r, userHandler, postHandler, tokenManager)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: r,
	}

	return s
}

func (s *Server) Start() error {
	s.log.Infof("starting HTTP server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := s.db.HealthCheck(ctx); err != nil {
		s.log.Errorf("healthcheck failed: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
