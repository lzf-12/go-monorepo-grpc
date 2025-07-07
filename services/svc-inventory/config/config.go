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
		Database  Database `json:"database"`
		Redis     Redis    `json:"redis"`
	}
	Database struct {
		InitSeeds bool   `json:"init_seeds"`
		Uri       string `json:"uri"`
	}
	Redis struct {
		Uri string `json:"uri"`
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
		Port:      env.Get("PORT", "50051").String(),

		Database: Database{
			InitSeeds: env.Get("INIT_SEEDS", "false").Bool(),
			Uri:       env.Get("DB_URI", "").String(),
		},
	}

	return cfg, nil
}
