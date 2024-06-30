package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	// Replace with your mkcert-generated certificate and key
	cert, err := tls.LoadX509KeyPair("localhost+2.pem", "localhost+2-key.pem")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3", "webtransport"},
	}

	h3Server := http3.Server{
		Addr:      ":4433",
		TLSConfig: tlsConfig,
	}

	wtServer := webtransport.Server{
		H3: h3Server,
	}

	http.HandleFunc("/webtransport", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received WebTransport request")
		session, err := wtServer.Upgrade(w, r)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebTransport", http.StatusInternalServerError)
			return
		}

		go handleSession(session)
	})

	log.Println("Starting WebTransport server on port 4433")
	if err := h3Server.ListenAndServeTLS("localhost+2.pem", "localhost+2-key.pem"); err != nil {
		log.Fatalf("Failed to start WebTransport server: %v", err)
	}
}

func handleSession(session *webtransport.Session) {
	for {
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println("Failed to accept stream:", err)
			return
		}
		go handleStream(stream)
	}
}

func handleStream(stream webtransport.Stream) {
	defer stream.Close()

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		log.Println("Failed to read from stream:", err)
		return
	}

	fmt.Printf("Received: %s\n", string(buf[:n]))

	response := fmt.Sprintf("Echo: %s", string(buf[:n]))
	_, err = stream.Write([]byte(response))
	if err != nil {
		log.Println("Failed to write to stream:", err)
	}
}
