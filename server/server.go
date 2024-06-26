package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webtransport", webTransportHandler)

	log.Println("Starting WebTransport server on port 4433")
	err := http3.ListenAndServeQUIC("0.0.0.0:4433", "cert.pem", "key.pem", mux)
	if err != nil {
		log.Fatalf("Failed to start HTTP/3 server: %v", err)
	}
}

func webTransportHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received WebTransport request")
	conn, ok := r.Context().Value(http3.ServerContextKey).(quic.Connection)
	if !ok {
		http.Error(w, "Not a WebTransport request", http.StatusBadRequest)
		return
	}

	go func() {
		for {
			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				log.Println(err)
				return
			}
			go handleStream(stream)
		}
	}()
}

func handleStream(stream quic.Stream) {
	defer stream.Close()

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		log.Println(err)
		return
	}

	fmt.Printf("Received: %s\n", string(buf[:n]))

	response := fmt.Sprintf("Echo: %s", string(buf[:n]))
	_, err = stream.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func generateTLSConfig() *tls.Config {
	certPath := "cert.pem"
	keyPath := "key.pem"

	log.Printf("Loading cert from: %s", certPath)
	log.Printf("Loading key from: %s", keyPath)

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Fatalf("Certificate file does not exist: %s", certPath)
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		log.Fatalf("Key file does not exist: %s", keyPath)
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load key pair: %s", err)
	}

	log.Println("Successfully loaded the key pair")

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"},
	}
}
