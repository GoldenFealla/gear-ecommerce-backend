package postgres

import "github.com/jackc/pgx/v5"

type GearRepository struct {
	Conn *pgx.Conn
}

func NewGearRepository(conn *pgx.Conn) *GearRepository {
	return &GearRepository{Conn: conn}
}
