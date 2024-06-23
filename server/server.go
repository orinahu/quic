package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"github.com/quic-go/quic-go"
)

func main() {
	listener, err := quic.ListenAddr("localhost:4433", generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("QUIC server listening on port 4433")

	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for {
				stream, err := session.AcceptStream(context.Background())
				if err != nil {
					log.Fatal(err)
				}

				go handleStream(stream)
			}
		}()
	}
}

func handleStream(stream quic.Stream) {
	defer stream.Close()

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	fmt.Printf("Received: %s\n", string(buf[:n]))

	response := fmt.Sprintf("Echo: %s", string(buf[:n]))
	_, err = stream.Write([]byte(response))
	if err != nil {
		log.Fatal(err)
	}
}

func generateTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
