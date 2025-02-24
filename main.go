package main

import (
	"log/slog"
	"os"
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

	Client *imapclient.Client
)

func main() {
	// Changes the logger to charmbracelet/log
	// Call slog so logging can be agnostic
	logger := log.NewWithOptions(os.Stderr, log.Options{
		TimeFormat:      time.DateTime,
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})

	slogger := slog.New(logger)
	slog.SetDefault(slogger)
	log.SetDefault(logger)

	// Load .env (for convenience)
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env Detected; continuing...")
	}

	// Parse Env variables
	IMAPUser = os.Getenv("IMAP_USER")
	IMAPPass = os.Getenv("IMAP_PASSWORD")
	IMAPAddress = os.Getenv("IMAP_ADDRESS")
	IMAPCertificate := os.Getenv("IMAP_CERTIFICATE")
	// Required
	if IMAPUser == "" {
		log.Fatal("`IMAP_USER` not set, required variable")
	}
	if IMAPPass == "" {
		log.Fatal("`IMAP_PASSWORD` not set, required variable")
	}
	// Optional
	if IMAPAddress == "" {
		log.Warn("No `IMAP_ADDRESS` set, using default 127.0.0.1:1143")
		IMAPAddress = "127.0.0.1:1143"
	}
	if IMAPCertificate == "" {
		log.Warn("No `IMAP_CERTIFICATE` set, using insecure TLS")
	}

	// Init Connection
	Client = connectIMAP(IMAPAddress, IMAPUser, IMAPPass, IMAPCertificate)
	defer Client.Close()

	// remove from memory
	IMAPPass = ""

	// Run main loop
	loop()
}

func loop() {
	messages := fetchMessages(Client)
	
	// Parse those messages
	for _, msg := range messages {
		log.Printf("Subject: %v", msg.Envelope.Subject)
		log.Printf("- Sender: %v", msg.Envelope.From)
	}
}
