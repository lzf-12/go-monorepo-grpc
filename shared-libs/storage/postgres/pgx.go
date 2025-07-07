package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PostgresPgx struct {
		pool *pgxpool.Pool
	}

	PgxRows = pgx.Rows
	PgxTx   = pgx.Tx
)

var PgxErrNoRows = pgx.ErrNoRows

const (
	defaultMaxOpen     int32 = 20
	defaultMinIdle     int32 = 10
	defaultMaxIdleTime       = 1 * time.Hour
)

// NewPgx creates a new pgx connection pool.
// maxopen, maxidle and maxlifetime are optional. When nil a sensible default is
// applied (20, 10 and 1h respectively).
func NewPgx(dsn string, maxopen, maxidle *int, maxlifetime *time.Duration) (*PostgresPgx, error) {
	if err := validateDSN(dsn); err != nil {
		return nil, err
	}

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Apply caller‑supplied pool sizing (or defaults).
	if maxopen != nil {
		poolCfg.MaxConns = int32(*maxopen)
	} else {
		poolCfg.MaxConns = defaultMaxOpen
	}

	if maxidle != nil {
		poolCfg.MinConns = int32(*maxidle)
	} else {
		poolCfg.MinConns = defaultMinIdle
	}

	if maxlifetime != nil {
		poolCfg.MaxConnIdleTime = *maxlifetime
	} else {
		poolCfg.MaxConnIdleTime = defaultMaxIdleTime
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	// verify the pool is reachable before returning.
	pg := &PostgresPgx{pool: pool}
	if err := pg.Ping(); err != nil {
		pool.Close()
		return nil, err
	}

	return pg, nil
}

// checks database liveness with a 5‑second timeout.
func (p *PostgresPgx) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.pool.Ping(ctx)
}

// gracefully shuts down the pool and waits for all
// in‑flight queries to finish.
func (p *PostgresPgx) Close() {
	p.pool.Close()
}

// to confirm the connection is usable.
func (p *PostgresPgx) IsReady() error {
	var dummy int
	return p.pool.QueryRow(context.Background(), "SELECT 1").Scan(&dummy)
}

// return *pgxpool.Pool
func (p *PostgresPgx) Pool() *pgxpool.Pool {
	return p.pool
}
