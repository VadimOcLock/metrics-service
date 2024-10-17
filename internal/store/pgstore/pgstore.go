package pgstore

import (
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewPgStore(db *pgxpool.Pool) metricservice.Store {
	return PgStore{
		Queries: New(db),
		db:      db,
	}
}

var _ metricservice.Store = (*PgStore)(nil)
