package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"vault/db"
)

const (
	encryptExt = "crypt"
)

var (
	// Options:- dev, prod
	macEnv = "dev"

	// Domains to whitelist
	domains = []string{"file.darrenyxj.com"}
	// domains = []string{"darrenyxj.com", "file.darrenyxj.com"}
)

func main() {
	fmt.Println("Vault Start")

	// Application setup
	appInit()

	// Initialize DB
	db.InitMainDb()
	defer db.CloseMainDB()

	// Define http routes
	http.HandleFunc("/encrypt", encryptFunc)
	http.HandleFunc("/decrypt", decryptFunc)
	http.HandleFunc("/getcounts", getBothCounts)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	// Listen and Serve
	if macEnv == "prod" {
		// Production mode with TLS
		fmt.Println("Production Mode:- TLS enabled")
		cm := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      autocert.DirCache("vault-autocert"),
		}

		server := &http.Server{
			Addr: ":443",
			TLSConfig: &tls.Config{
				GetCertificate: cm.GetCertificate,
			},
		}

		go func() {
			log.Fatal(http.ListenAndServe(":80", cm.HTTPHandler(nil)))
		}()
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else if macEnv == "dev" {
		// Development mode without TLS
		fmt.Println("Development Mode:- TLS disabled")
		log.Fatal(http.ListenAndServe(":80", nil))
	} else {
		log.Fatalf("Invalid environment variable value VAULT_MAC_ENV %v. Leave it unset or choose between dev or prod.", macEnv)
	}
}

func appInit() {
	// Read environment variables
	evMacEnv := os.Getenv("VAULT_MAC_ENV")
	if evMacEnv != "" {
		macEnv = evMacEnv
	}
}
