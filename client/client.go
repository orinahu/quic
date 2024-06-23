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
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	// Create a QUIC session to the server
	session, err := quic.DialAddr(context.Background(), "localhost:4433", tlsConfig, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Open a stream within the QUIC session
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	message := "Hello, QUIC!"
	fmt.Printf("Sending: %s\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	fmt.Printf("Received: %s\n", string(buf[:n]))
}
