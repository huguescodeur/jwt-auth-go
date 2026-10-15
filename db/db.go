package db

import (
	"database/sql"
	"fmt"
)

func NewSQLStorage(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("error ping: %w", err)
	}

	fmt.Println("Database connection established")

	return db, nil
}
