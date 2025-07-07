package config

import (
	"log"
	"ops-monorepo/shared-libs/env"
)

type (
	Config struct {
		Env       string   `json:"env"`
		AppName   string   `json:"app_name"`
		DebugMode bool     `json:"debug_mode"`
		Port      string   `json:"port"`
		HTTPPort  string   `json:"http_port"`
		Database  Database `json:"database"`
		JWT       JWT      `json:"jwt"`
	}
	Database struct {
		InitSeeds bool   `json:"init_seeds"`
		Uri       string `json:"uri"`
	}
	JWT struct {
		Secret               string `json:"secret"`
		AccessTokenDuration  int    `json:"access_token_duration"`  // in minutes
		RefreshTokenDuration int    `json:"refresh_token_duration"` // in hours
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
		Port:      env.Get("PORT", "50053").String(),
		HTTPPort:  env.Get("HTTP_PORT", "8080").String(),

		Database: Database{
			InitSeeds: env.Get("INIT_SEEDS", "false").Bool(),
			Uri:       env.Get("DB_URI", "").String(),
		},

		JWT: JWT{
			Secret:               env.Get("JWT_SECRET", "your-secret-key").String(),
			AccessTokenDuration:  int(env.Get("JWT_ACCESS_DURATION", "15m").DurationInSecond()),
			RefreshTokenDuration: int(env.Get("JWT_REFRESH_DURATION", "24h").DurationInSecond()),
		},
	}

	return cfg, nil
}
