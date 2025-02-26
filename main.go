package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/joho/godotenv"
)

var (
	IMAPUser        string // IMAP_USER
	IMAPPass        string // IMAP_PASSWORD
	IMAPAddress     string // IMAP_ADDRESS, default 127.0.0.1:1443
	IMAPCertificate string // IMAP_CERTIFICATE

	DB_URL string // DB_URL
)

func main() {
	// Init Logging
	initLogging()

	// Init Variables
	parseENV()
	parseCSV()

	// Init Database
	initDB(DB_URL)

	// Init Connection
	client := connectIMAP(IMAPAddress, IMAPUser, IMAPPass, IMAPCertificate)
	defer client.Close()

	// Run main loop
	loop(client)

	// Cleanup
	client.Logout().Wait()
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

func loop(c *imapclient.Client) {
	messages := fetchMessages(c)

	for _, msg := range messages {
		// admissions@osu.edu -> osu.edu
		address := strings.Split(msg.Envelope.From[0].Addr(), "@")

		// osu.edu -> CollegeRecord
		record := fetchRecord(address[len(address)-1])
		log.Debugf("- %v", record.name)

		// Mark as Seen (signal as already processed)
		/*
			if err := markSeen(c, msg); err != nil {
				log.Fatalf("Failed to mark message as seen: %v", err)
			}
		*/
	}
}
