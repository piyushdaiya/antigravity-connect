package main

import (
	"log"
	"github.com/piyushdaiya/antigravity-connect/internal/certs"
	"github.com/piyushdaiya/antigravity-connect/internal/server"
	"net/http"
)

func main() {
	// 1. Generate Certs
	tlsConfig, err := certs.GenerateTLSConfig()
	if err != nil {
		log.Fatalf("Failed to generate certs: %v", err)
	}

	// 2. Prepare Server
	srv := &http.Server{
		TLSConfig: tlsConfig,
	}

	// 3. Start Application
	// Make sure your HTML files are in the /web folder before building!
	server.Start("3000", srv)
}
