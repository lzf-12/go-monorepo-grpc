package internal

import (
	"ops-monorepo/shared-libs/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() {

	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())

	// JWT Authentication middleware configuration
	authConfig := middleware.AuthConfig{
		UserServiceURL: "svc-user:50053", // Use service name for Docker, or localhost for local dev
		Timeout:        5 * time.Second,
	}

	api := s.Router.Group("api")
	v1 := api.Group("v1")

	// Protected routes that require authentication
	protected := v1.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(authConfig))
	{
		// Create order endpoint requires authentication
		protected.POST("/orders", s.order.handler.CreateOrder)

		// You can add role-based protection like this:
		// protected.POST("/orders", middleware.RequireRole("user", "admin"), s.order.handler.CreateOrder)
	}

	// Public routes (no authentication required)
	// v1.GET("/health", s.HealthCheck)
}
