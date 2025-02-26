package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

// DSN is the connection string.
// Example: postgres://postgres:@localhost:5432/test?sslmode=disable
func initDB(dsn string) {
	// Create Connection
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	// Table SQL
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS colleges (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		domain TEXT NOT NULL,
		state_abbr CHAR(2) NOT NULL REFERENCES states(abbreviation) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS states (
		abbreviation CHAR(2) PRIMARY KEY,
		name TEXT UNIQUE NOT NULL
	);
	`

	// Run Migration
	_, err = db.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Index SQL
	createIndexQuery := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_colleges_domain ON colleges(domain);
	`
	// Create Index
	_, err = db.Exec(context.Background(), createIndexQuery)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}
