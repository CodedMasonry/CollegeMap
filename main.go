package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

var (
	IMAPUser    string // IMAP_USER
	IMAPPass    string // IMAP_PASSWORD
	IMAPAddress string // IMAP_ADDRESS, default 0.0.0.0:1443
)

func main() {
	// Changes the logger to charmbracelet/log
	// Call slog so logging can be agnostic
	logger := log.NewWithOptions(os.Stderr, log.Options{
		TimeFormat:      time.DateTime,
		ReportTimestamp: true,
	})
	slogger := slog.New(logger)
	slog.SetDefault(slogger)

	// Load .env (for convenience)
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env Detected; continuing...")
	}

	// Parse Env variables
	IMAPUser = os.Getenv("IMAP_USER")
	IMAPPass = os.Getenv("IMAP_PASSWORD")
	IMAPAddress = os.Getenv("IMAP_ADDRESS")
	if IMAPUser == "" {
		log.Fatal("`IMAP_USER` not set, required variable")
	}
	if IMAPPass == "" {
		log.Fatal("`IMAP_PASSWORD` not set, required variable")
	}
	if IMAPAddress == "" {
		log.Warn("No `IMAP_ADDRESS` set, using default 0.0.0.0:1143")
		IMAPAddress = "0.0.0.0:1143"
	}
}
