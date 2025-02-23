package main

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/charmbracelet/log"
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
