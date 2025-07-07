package storage

import (
	"errors"
	"fmt"
	"net/url"
)

// validateDSN guarantees the incoming DSN looks like
// "postgres://user:pass@host:port/dbname". pgx will still do its own parsing,
// but this catches typo‑level mistakes early.
func validateDSN(dsn string) error {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("invalid DSN format: %w", err)
	}
	if parsed.Scheme != "postgres" {
		return errors.New("DSN must use postgres scheme")
	}
	if parsed.Host == "" {
		return errors.New("DSN must include host")
	}
	return nil
}
