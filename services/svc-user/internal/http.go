package internal

import (
	"ops-monorepo/services/svc-user/config"
	"ops-monorepo/services/svc-user/internal/delivery/http"
	"ops-monorepo/shared-libs/logger"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Router *mux.Router
	user   *userImpl
	Log    logger.Logger
}

func NewHTTPServer(cfg *config.Config) *httpServer {

	dep := InitDependencies(cfg)

	r := mux.NewRouter()

	return &httpServer{
		Router: r,
		user:   &dep.Impl.userImpl,
		Log:    dep.log,
	}
}

func (s *httpServer) Register() {

	// user HTTP endpoints
	userHandler := http.NewUserHTTPHandler(s.user.usecase, s.Log)
	userHandler.RegisterRoutes(s.Router)
}