package postgresDB

import (
	"fmt"
	"kaf-interface/internal/orders/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.Config) (*sqlx.DB, error) {

	loginString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	db, err := sqlx.Open("postgres", loginString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migration(db); err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return db, nil
}
