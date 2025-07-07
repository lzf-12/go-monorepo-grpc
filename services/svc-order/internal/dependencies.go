package internal

import (
	"errlib"
	"log"
	"ops-monorepo/services/svc-order/config"
	"ops-monorepo/services/svc-order/internal/delivery/handler"
	"ops-monorepo/services/svc-order/internal/repository"
	"ops-monorepo/services/svc-order/internal/usecase"
	"ops-monorepo/services/svc-order/seeds"
	"ops-monorepo/services/svc-order/validator"
	"ops-monorepo/shared-libs/jwt"
	"ops-monorepo/shared-libs/logger"
	"os"
	inventoryv1 "pb_schemas/inventory/v1"
	userv1 "pb_schemas/user/v1"

	gg "ops-monorepo/shared-libs/grpc/client"

	"time"

	pg "ops-monorepo/shared-libs/storage/postgres"
)

type Dependencies struct {
	DbPostgres   *pg.Postgres
	ErrorHandler *errlib.ErrorHandler
	Jwt          *jwt.TokenManager
	GrpcDeps     GrpcDeps
	Impl         Impl
}

type GrpcDeps struct {
	InventoryGrpcClient inventoryv1.InventoryServiceClient
	UserGrpcClient      userv1.UserServiceClient
}

type Impl struct {
	Order
}

type Order struct {
	handler    handler.IOrder
	usecase    usecase.IOrderUsecase
	repository repository.IOrderSQLRepository
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

	// errorHandler
	dep.ErrorHandler = errlib.NewErrorHandler(cfg.DebugMode)
	zl.Info("error handler ok..")

	// seeds db when enabled
	if cfg.Database.InitSeeds {
		zl.Info("database seeds flag enabled, executing seeds scripts... ")
		if err := seeds.ExecuteDefaultTableScripts(db, zl); err != nil {
			zl.Errorf("error when executing db seeds, %v", err)
		}

	}

	// grpc services
	clientRegistry := gg.NewClientRegistry()
	invConn, err := clientRegistry.GetConnection(cfg.GrpcServices.ServiceInventoryGrpcUrl)
	if err != nil || invConn == nil {
		zl.Fatal("cannot establish connection with inventory service..")
	}
	inventoryClient := inventoryv1.NewInventoryServiceClient(invConn)
	dep.GrpcDeps.InventoryGrpcClient = inventoryClient
	zl.Info("inventory grpc client ok..")

	// validator
	val := validator.NewValidator()

	//order
	dep.Impl.Order.repository = repository.NewOrderRepository(db)
	dep.Impl.Order.usecase = usecase.NewOrderUsecase(dep.Impl.Order.repository, zl, dep.GrpcDeps.InventoryGrpcClient)
	dep.Impl.Order.handler = handler.NewOrderHandler(val, zl, dep.ErrorHandler, dep.Impl.usecase)
	zl.Info("order module ok..")

	return dep
}
