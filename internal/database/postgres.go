package database

import (
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/config"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	Sql *sql.DB
}

var (
	dbInstance *Postgres
)

func NewPostgres(cfg *config.Config) *Postgres {
	username := cfg.DB_USERNAME
	password := cfg.DB_PASSWORD
	host := cfg.DB_HOST
	port := cfg.DB_PORT
	database := cfg.DB_DATABASE
	schema := cfg.DB_SCHEMA
	// Reuse Connection
	if dbInstance != nil {
		log.Println("Reusing database connection")
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &Postgres{
		Sql: db,
	}
	return dbInstance
}
