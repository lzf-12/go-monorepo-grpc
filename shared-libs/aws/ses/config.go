package aws_ses

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type SmtpConfig struct {
	Host, Username, Password string
	Port                     int           // 587 (STARTTLS) or 465 (TLS wrapper)
	MaxConns                 int           // upper bound for live TCP conns
	ConnectTimeout           time.Duration // SMTP dial timeout
	SendTimeout              time.Duration // per-message timeout
}

func LoadEnvSmtpConfig() (*SmtpConfig, error) {
	config := &SmtpConfig{}

	config.Host = os.Getenv("SMTP_HOST")
	config.Username = os.Getenv("SMTP_USERNAME")
	config.Password = os.Getenv("SMTP_PASSWORD")

	// port with default fallback
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
		}
		config.Port = port
	} else {
		config.Port = 587 // default STARTTLS port
	}

	// MaxConns with default fallback
	if maxConnsStr := os.Getenv("SMTP_MAX_CONNS"); maxConnsStr != "" {
		maxConns, err := strconv.Atoi(maxConnsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SMTP_MAX_CONNS: %w", err)
		}
		config.MaxConns = maxConns
	} else {
		config.MaxConns = 10 // default max connections
	}

	// ConnectTimeout with default fallback
	if timeoutStr := os.Getenv("SMTP_CONNECT_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SMTP_CONNECT_TIMEOUT: %w", err)
		}
		config.ConnectTimeout = timeout
	} else {
		config.ConnectTimeout = 30 * time.Second // default 30 seconds
	}

	// SendTimeout with default fallback
	if timeoutStr := os.Getenv("SMTP_SEND_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SMTP_SEND_TIMEOUT: %w", err)
		}
		config.SendTimeout = timeout
	} else {
		config.SendTimeout = 60 * time.Second // default 60 seconds
	}

	return config, nil
}
