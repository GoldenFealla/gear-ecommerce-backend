package postgres

import "github.com/jackc/pgx/v5"

type UserRepository struct {
	Conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{Conn: conn}
}
