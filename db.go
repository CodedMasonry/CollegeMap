package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

// Actual Database connection

// DSN is the connection string.
// Example: postgres://postgres:@localhost:5432/test?sslmode=disable
func initDB(dsn string) (db *pgx.Conn) {
	// Create Connection
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	// States
	// - abbreviation: shorthand for state
	// - name:         Name of the state
	createStatesTableQuery := `
	CREATE TABLE IF NOT EXISTS states (
		abbreviation CHAR(2) PRIMARY KEY,
		name VARCHAR(32) UNIQUE NOT NULL
	);
	`

	// Colleges
	// - id:         random int
	// - name:       Name of College
	// - domain:     domain name, ex: osu.edu
	// - is_ivory:   whether it is an ivory school
	// - num_emails: the number of emails
	// - state_abbr: state associated with school
	createCollegeTableQuery := `
	CREATE TABLE IF NOT EXISTS colleges (
		id SERIAL PRIMARY KEY,
		name VARCHAR(128) NOT NULL,
		domain VARCHAR(32) NOT NULL,
		is_ivory BOOLEAN NOT NULL,
		num_emails INTEGER NOT NULL,
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
	return
}

// Increment if record exists, else create new record
func incrementCollegeEmails(conn *pgx.Conn, rec *CollegeRecord) error {
	checkStateExists(conn, rec.stateAbbr)

	query := `
		INSERT INTO colleges (name, domain, is_ivory, num_emails, state_abbr)
		VALUES ($1, $2, $3, 1, $4)
		ON CONFLICT (domain)
		DO UPDATE SET num_emails = colleges.num_emails + 1;
	`
	_, err := conn.Exec(context.Background(), query, rec.name, rec.domain, rec.isIvory, rec.stateAbbr)
	return err
}

// Checks whether the state exists, if not create a record
func checkStateExists(conn *pgx.Conn, abbreviation string) error {
	stateName := states[abbreviation]

	query := `
		INSERT INTO states (abbreviation, name)
		VALUES ($1, $2)
		ON CONFLICT (abbreviation) DO NOTHING;
	`
	_, err := conn.Exec(context.Background(), query, abbreviation, stateName)
	return err
}