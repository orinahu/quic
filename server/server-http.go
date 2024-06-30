package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	certFile := "cert.pem"
	keyFile := "key.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load key pair: %s", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      ":4433",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/your-endpoint", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		var requestData struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		responseData := map[string]string{
			"response": "Received: " + requestData.Message,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	})

	log.Println("Starting HTTPS server on port 4433")
	log.Fatal(server.ListenAndServeTLS("localhost.crt", "localhost.decrypted.key"))
}
