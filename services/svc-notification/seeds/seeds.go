package seeds

import (
	"context"
	"fmt"
	"ops-monorepo/shared-libs/logger"
	pg "ops-monorepo/shared-libs/storage/postgres"
	"os"
	"path/filepath"
	"runtime"
)

func ExecuteDefaultTableScripts(db *pg.PostgresPgx, log logger.Logger) error {
	// Get the current file's directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}
	dir := filepath.Dir(filename)

	// Read and execute default_tables.sql
	tablesPath := filepath.Join(dir, "default_tables.sql")
	tablesSQL, err := os.ReadFile(tablesPath)
	if err != nil {
		return fmt.Errorf("failed to read default_tables.sql: %w", err)
	}

	if _, err := db.Pool().Exec(context.Background(), string(tablesSQL)); err != nil {
		return fmt.Errorf("failed to execute default_tables.sql: %w", err)
	}

	log.Info("default_tables.sql executed successfully")

	// Read and execute default_data.sql
	dataPath := filepath.Join(dir, "default_data.sql")
	dataSQL, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("failed to read default_data.sql: %w", err)
	}

	if _, err := db.Pool().Exec(context.Background(), string(dataSQL)); err != nil {
		return fmt.Errorf("failed to execute default_data.sql: %w", err)
	}

	log.Info("default_data.sql executed successfully")

	return nil
}