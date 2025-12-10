package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"service/internal/configs"
	"time"

	"github.com/pressly/goose/v3"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

func InitReformDB(cnf configs.Database) (*reform.DB, error) {
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.User,
		cnf.Password,
		cnf.Host,
		cnf.Port,
		cnf.Name,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(cnf.MaxOpenConnection)
	db.SetConnMaxLifetime(time.Minute * time.Duration(cnf.MaxLifeTime))

	if err = initMigration(db); err != nil {
		return nil, err
	}

	//про логгер подумать
	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	reformDB := reform.NewDB(db, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))

	return reformDB, nil
}

func initMigration(db *sql.DB) error {
	err := goose.Up(db, "migrations")
	if err != nil {
		return fmt.Errorf("failed up migrations: %w", err)
	}
	return nil
}
