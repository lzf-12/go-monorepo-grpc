package config

import (
	"log"
	"ops-monorepo/shared-libs/env"
)

type (
	Config struct {
		Env          string       `json:"env"`
		AppName      string       `json:"app_name"`
		DebugMode    bool         `json:"debug_mode"`
		Port         string       `json:"port"`
		Database     Database     `json:"database"`
		Redis        Redis        `json:"redis"`
		GrpcServices GrpcServices `json:"grpc_services"`
	}
	Database struct {
		InitSeeds bool   `json:"init_seeds"`
		Uri       string `json:"uri"`
	}
	Redis struct {
		Uri string `json:"uri"`
	}

	GrpcServices struct {
		ServiceUserGrpcUrl         string `json:"service_user_grpc_url"`
		ServiceInventoryGrpcUrl    string `json:"service_inventory_grpc_url"`
		ServiceNotificationGrpcUrl string `json:"service_notification_grpc_url"`
	}
)

func LoadConfig(path string) (*Config, error) {

	if err := env.LoadEnv(path); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := &Config{
		Env:       env.Get("ENV", "").String(),
		AppName:   env.Get("APPNAME", "").String(),
		DebugMode: env.Get("DEBUG_MODE", "").Bool(),
		Port:      env.Get("PORT", "8080").String(),

		Database: Database{
			InitSeeds: env.Get("INIT_SEEDS", "false").Bool(),
			Uri:       env.Get("DB_URI", "").String(),
		},

		Redis: Redis{
			Uri: env.Get("REDIS_URI", "").String(),
		},

		GrpcServices: GrpcServices{
			ServiceUserGrpcUrl:         env.Get("SERVICE_USER_GRPC_URL", "").String(),
			ServiceInventoryGrpcUrl:    env.Get("SERVICE_INVENTORY_GRPC_URL", "").String(),
			ServiceNotificationGrpcUrl: env.Get("SERVICE_NOTIFICATION_GRPC_URL", "").String(),
		},
	}

	return cfg, nil
}
