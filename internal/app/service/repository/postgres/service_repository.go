package postgres

import (
	"github.com/jackc/pgx"

	"github.com/alSergey/TechMain_2021_db_forum/internal/app/service"
)

type ServiceRepository struct {
	conn *pgx.ConnPool
}

func NewServiceRepository(conn *pgx.ConnPool) service.ServiceRepository {
	return &ServiceRepository{
		conn: conn,
	}
}
