package database

import "github.com/jackc/pgx/v4"

type DB struct {
	Connection *pgx.Conn
}
