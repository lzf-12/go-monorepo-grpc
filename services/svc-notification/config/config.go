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
		SMTP      SMTP     `json:"smtp"`
	}
	Database struct {
		InitSeeds bool   `json:"init_seeds"`
		Uri       string `json:"uri"`
	}
	SMTP struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		From     string `json:"from"`
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
		Port:      env.Get("PORT", "50052").String(),

		Database: Database{
			InitSeeds: env.Get("INIT_SEEDS", "false").Bool(),
			Uri:       env.Get("DB_URI", "").String(),
		},

		SMTP: SMTP{
			Host:     env.Get("SMTP_HOST", "").String(),
			Port:     env.Get("SMTP_PORT", "587").IntDefault(587),
			Username: env.Get("SMTP_USERNAME", "").String(),
			Password: env.Get("SMTP_PASSWORD", "").String(),
			From:     env.Get("SMTP_FROM", "").String(),
		},
	}

	return cfg, nil
}
