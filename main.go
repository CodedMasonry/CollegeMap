package main

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
	db := initDB(DB_URL)

	// Init Connection
	client := connectIMAP(IMAPAddress, IMAPUser, IMAPPass, IMAPCertificate)
	defer client.Close()

	// Run main loop
	loop(client, db)

	// Cleanup
	client.Logout().Wait()
}
