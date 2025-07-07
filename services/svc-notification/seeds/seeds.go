package seeds

import (
	"context"
	_ "embed"
	"ops-monorepo/shared-libs/logger"
	pg "ops-monorepo/shared-libs/storage/postgres"
	storage "ops-monorepo/shared-libs/storage/postgres"
)

//go:embed default_tables.sql
var DefaultTableSQL []byte

//go:embed default_data.sql
var DefaultDataSQL []byte

var CheckIfSchemaExist string

func ExecuteDefaultTableScripts(db *pg.PostgresPgx, log logger.Logger) error {

	// check if schemas already exist
	CheckIfSchemaExist = `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'order_service';`

	var rows storage.PgxRows

	rows, err := db.Pool().Query(context.Background(), CheckIfSchemaExist)
	if err != nil {
		return err
	}
	defer rows.Close()

	schemaTables := []string{}
	for rows.Next() {
		var table string
		rows.Scan(&table)
		schemaTables = append(schemaTables, table)
	}

	// if schema already has table, skip
	if len(schemaTables) > 0 {
		log.Info("schema and tables already exist, seeds script skipped..")
		return nil
	}

	if len(DefaultTableSQL) < 1 {
		log.Info("seeds sql not loaded")
	}

	// proceed if schema and tables not exist
	_, err = db.Pool().Exec(context.Background(), string(DefaultTableSQL))
	if err != nil {
		log.Fatalf("create default table failed: %v", err)
		return err
	}

	_, err = db.Pool().Exec(context.Background(), string(DefaultDataSQL))
	if err != nil {
		log.Fatalf("seeding default data failed: %v", err)
		return err
	}

	log.Info("seeds scripts successfully executed..")

	return nil
}
