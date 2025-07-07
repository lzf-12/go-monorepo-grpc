package internal

import (
	"log"

	"ops-monorepo/services/svc-user/config"
	"ops-monorepo/services/svc-user/internal/delivery/handler"
	"ops-monorepo/services/svc-user/internal/repository"
	"ops-monorepo/services/svc-user/internal/usecase"
	"ops-monorepo/services/svc-user/seeds"
	grpcErr "ops-monorepo/shared-libs/grpc/errors"
	"ops-monorepo/shared-libs/logger"
	"os"
	"time"

	pg "ops-monorepo/shared-libs/storage/postgres"
)

type Dependencies struct {
	log            logger.Logger
	DbPostgres     *pg.PostgresPgx
	GrpcErrHandler *grpcErr.GRPCErrorHandler
	Impl           Impl
}

type Impl struct {
	userImpl
}

type userImpl struct {
	handler    handler.IUserHandler
	usecase    usecase.IUserUsecase
	repository repository.IUserSQLRepository
}

func InitDependencies(cfg *config.Config) Dependencies {

	if cfg == nil {
		log.Fatalf("error configuration is nil")
	}

	dep := Dependencies{}

	// logger
	logCfg := &logger.Config{
		Level:       "info",
		Format:      "console",
		Output:      os.Stdout,
		TimeFormat:  time.RFC3339,
		Caller:      false,
		ServiceName: cfg.AppName,
		Version:     "",
	}

	if cfg.DebugMode {
		logCfg.Level = "debug"
	}

	zl := logger.New(logCfg)
	dep.log = zl
	zl.Info("logger ok..")

	// db pool connection
	db, err := pg.NewPgx(
		cfg.Database.Uri,
		nil, nil, nil, // use default connection limiter
	)
	if err != nil || db == nil {
		zl.Fatal("error failed to initialize postgres connection")
	}
	zl.Info("postgres ok..")

	// grpcErrHandler
	dep.GrpcErrHandler = grpcErr.NewGRPCErrorHandler()
	zl.Info("grpcErr handler ok..")

	// seeds db when enabled
	if cfg.Database.InitSeeds {
		zl.Info("database seeds flag enabled, executing seeds scripts... ")
		if err := seeds.ExecuteDefaultTableScripts(db, zl); err != nil {
			zl.Errorf("error when executing db seeds, %v", err)
		}
	}
	zl.Info("seeds ok..")

	// user
	dep.Impl.userImpl.repository = repository.NewUserRepository(db)
	dep.Impl.userImpl.usecase = usecase.NewUserUsecase(zl, dep.Impl.repository, cfg)
	dep.Impl.userImpl.handler = handler.NewUserHandler(zl, dep.Impl.usecase, dep.GrpcErrHandler)
	zl.Info("user ok..")

	return dep
}