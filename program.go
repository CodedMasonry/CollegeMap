package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// Main func loop
func run(c *imapclient.Client, db *pgx.Conn) {
	messages := fetchMessages(c)
	numMessages := len(messages)

	for _, msg := range messages {
		// admissions@osu.edu -> admissions.osu.edu
		senderSplit := strings.Split(msg.Envelope.From[0].Addr(), "@")
		address := senderSplit[len(senderSplit)-1]

		// admissions.osu.edu -> [admissions osu edu]
		parts := strings.Split(address, ".")
		// [admissions osu edu] -> osu.edu
		domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

		// osu.edu -> CollegeRecord
		record := fetchRecord(domain)

		// If the domain doesn't actually exist, skip
		if record == nil {
			log.Warnf("! %v", domain)
			continue
		}
		log.Debugf("- %v", record.name)

		// CollegeRecord -> DB Entry
		err := incrementCollegeEmails(db, record)
		if err != nil {
			log.Fatalf("failed to increment record: %v", err)
		}

		// Mark as Seen (signal as already processed)
		if err := markSeen(c, msg); err != nil {
			log.Fatalf("Failed to mark message as seen: %v", err)
		}
	}
	log.Infof("%d Records Incremented", numMessages)
}

// Changes the logger to charmbracelet/log
func initLogging() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		TimeFormat:      time.DateTime,
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})

	// Call slog so logging can be agnostic
	slogger := slog.New(logger)
	slog.SetDefault(slogger)
	log.SetDefault(logger)
}

func parseENV() {
	// Load .env (for convenience)
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env Detected; continuing...")
	}

	// Parse Env variables
	IMAPUser = os.Getenv("IMAP_USER")
	IMAPPass = os.Getenv("IMAP_PASSWORD")
	IMAPAddress = os.Getenv("IMAP_ADDRESS")
	IMAPCertificate = os.Getenv("IMAP_CERTIFICATE")
	DB_URL = os.Getenv("DB_URL")

	// Required IMAP
	if IMAPUser == "" {
		log.Fatal("`IMAP_USER` not set, required variable")
	}
	if IMAPPass == "" {
		log.Fatal("`IMAP_PASSWORD` not set, required variable")
	}

	// Optional IMAP
	if IMAPAddress == "" {
		log.Warn("No `IMAP_ADDRESS` set, using default 127.0.0.1:1143")
		IMAPAddress = "127.0.0.1:1143"
	}
	if IMAPCertificate == "" {
		log.Warn("No `IMAP_CERTIFICATE` set, using insecure TLS")
	}

	// Database
	if DB_URL == "" {
		log.Fatal("No `DB_URL` set, requires postgres DSN")
	}
}
