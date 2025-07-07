package internal

import (
	"ops-monorepo/services/svc-order/config"

	"github.com/gin-gonic/gin"
)

// represents the API server with all dependencies
type Server struct {
	Config *config.Config
	Router *gin.Engine
	order  *Order
}

// creates a new server instance
func NewHTTPServer(cfg *config.Config) *Server {

	dep := InitDependencies(cfg)
	return &Server{
		Config: cfg,
		Router: gin.New(),
		order:  &dep.Impl.Order,
	}
}

func (s *Server) Start(addr string) error {
	s.Router.SetTrustedProxies(nil)
	return s.Router.Run(addr)
}
