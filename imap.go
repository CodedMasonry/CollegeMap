package main

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/charmbracelet/log"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

func connectIMAP(addr, user, pass, cert string) (client *imapclient.Client) {
	// Set TLS options
	options := &imapclient.Options{
		TLSConfig: &tls.Config{},
	}

	// Whether TLS is secure or not
	if cert != "" {
		pool := x509.NewCertPool()
		if success := pool.AppendCertsFromPEM([]byte(cert)); !success {
			log.Fatal("Failed to decode PEM certificate")
		}

		options.TLSConfig.RootCAs = pool
	} else {
		options.TLSConfig.InsecureSkipVerify = true
	}

	// Connect
	client, err := imapclient.DialStartTLS(addr, options)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	// Login
	if err := client.Login(user, pass).Wait(); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	return
}

func fetchMessages(c *imapclient.Client) (messages []*imapclient.FetchMessageBuffer) {
	// Select Mailbox
	selectedMbox, err := c.Select("All Mail", nil).Wait()
	if err != nil {
		log.Fatalf("failed to select 'All Mail' : %v", err)
	}
	log.Infof("%d messages", selectedMbox.NumMessages)

	// Search for messages from colleges
	// Check if they haven't been "seen" (processed) yet
	search := &imap.SearchCriteria{Body: []string{".edu"}, Not: []imap.SearchCriteria{{Flag: []imap.Flag{"\\Seen"}}}}

	// Run Search
	ids, err := c.UIDSearch(search, nil).Wait()
	if err != nil {
		log.Fatalf("Failed to find messages: %v", err)
	}

	// Make sure there are messages
	numMessages := len(ids.AllUIDs())
	log.Infof("%d college messages", numMessages)
	// No messages cause EOF error
	if numMessages == 0 {
		return
	}

	// Get messages
	messages, err = c.Fetch(ids.All, &imap.FetchOptions{Envelope: true}).Collect()
	if err != nil {
		log.Fatalf("Failed to fetch messages: %v", err)
	}

	return
}

func markSeen(c *imapclient.Client, msg *imapclient.FetchMessageBuffer) error {
	storeFlags := imap.StoreFlags{
		Op:    imap.StoreFlagsAdd,
		Flags: []imap.Flag{imap.FlagSeen},
	}

	return c.Store(imap.UIDSetNum(msg.UID), &storeFlags, nil).Close()
}
