package internal

import "github.com/gin-gonic/gin"

func (s *Server) RegisterRoutes() {

	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	api := s.Router.Group("api")
	v1 := api.Group("v1")

	v1.POST("/orders", s.order.handler.CreateOrder)
}
