package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/quic-go/quic-go"
)

func main() {
	// Create TLS configuration (skip certificate verification for local testing)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Only for testing! Don't use in production
		NextProtos:         []string{"quic-example"},
	}

	// Configure QUIC
	quicConfig := &quic.Config{}

	// Connect to QUIC server
	ctx := context.Background()
	conn, err := quic.DialAddr(ctx, "localhost:8082", tlsConfig, quicConfig)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	fmt.Println("Connected to QUIC server at localhost:8082")

	// Open a new stream
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		fmt.Printf("Failed to open stream: %v\n", err)
		return
	}
	defer stream.Close()

	// Start a goroutine to read from the server
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stream.Read(buffer)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Server closed the stream")
				} else {
					fmt.Printf("Error reading from server: %v\n", err)
				}
				os.Exit(0)
				return
			}
			fmt.Printf("Server: %s\n", string(buffer[:n]))
		}
	}()

	// Read user input and send to server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages to send (press Ctrl+C to exit):")
	for scanner.Scan() {
		message := scanner.Text()

		_, err := stream.Write([]byte(message))
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
