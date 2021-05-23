package database

import (
	"github.com/jackc/pgx"
)

type Postgres struct {
	postgresConnPool *pgx.ConnPool
}

func NewPostgres(dataSourceName string) (*Postgres, error) {
	pgxConnConfig, err := pgx.ParseConnectionString(dataSourceName)
	if err != nil {
		return nil, err
	}

	pgxConnConfig.PreferSimpleProtocol = true

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgxConnConfig,
		MaxConnections: 200,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		postgresConnPool: pool,
	}, nil
}

func (p *Postgres) GetDatabase() *pgx.ConnPool {
	return p.postgresConnPool
}

func (p *Postgres) Close() {
	p.postgresConnPool.Close()
}
