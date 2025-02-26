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
	createStatesTableQuery := `
	CREATE TABLE IF NOT EXISTS states (
		abbreviation CHAR(2) PRIMARY KEY,
		name TEXT UNIQUE NOT NULL
	);
	`

	createCollegeTableQuery := `
	CREATE TABLE IF NOT EXISTS colleges (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		domain TEXT NOT NULL,
		state_abbr CHAR(2) NOT NULL REFERENCES states(abbreviation) ON DELETE CASCADE
	);
	`

	// Run Migration
	log.Info("Creating Database Tables If Non-Existent...")
	_, err = db.Exec(context.Background(), createStatesTableQuery)
	if err != nil {
		log.Fatalf("Failed to create state table: %v", err)
	}
	_, err = db.Exec(context.Background(), createCollegeTableQuery)
	if err != nil {
		log.Fatalf("Failed to create college table: %v", err)
	}

	// Index SQL
	createIndexQuery := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_colleges_domain ON colleges(domain);
	`
	// Create Index
	log.Info("Creating Indexes If Non-Existent...")
	_, err = db.Exec(context.Background(), createIndexQuery)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	log.Info("Database Initialized")
}

func incrementRecord(rec CollegeRecord) {

}
