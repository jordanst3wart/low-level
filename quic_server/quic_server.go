package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/quic-go/quic-go"
)

// generateTLSConfig generates a self-signed certificate for testing purposes
func generateTLSConfig() (*tls.Config, error) {
	// Generate a private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 180), // Valid for 180 days
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}

	// Create a self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	// Encode the certificate and private key to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	// Load the certificate
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	// Create TLS configuration
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-example"},
	}, nil
}

func handleStream(stream quic.Stream) {
	defer stream.Close()

	fmt.Printf("New stream established: %d\n", stream.StreamID())

	// Send welcome message
	welcomeMsg := "Welcome to the QUIC server! Type something to get an echo response.\n"
	_, err := stream.Write([]byte(welcomeMsg))
	if err != nil {
		fmt.Printf("Error sending welcome message: %v\n", err)
		return
	}

	// Read from the stream
	buffer := make([]byte, 1024)
	for {
		n, err := stream.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Stream %d closed by client\n", stream.StreamID())
			} else {
				fmt.Printf("Error reading from stream: %v\n", err)
			}
			return
		}

		message := string(buffer[:n])
		fmt.Printf("Received from stream %d: %s\n", stream.StreamID(), message)

		// Echo the message back
		response := fmt.Sprintf("Echo: %s", message)
		_, err = stream.Write([]byte(response))
		if err != nil {
			fmt.Printf("Error sending response: %v\n", err)
			return
		}
	}
}

func main() {
	// Generate TLS configuration with self-signed certificate
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		fmt.Printf("Failed to generate TLS config: %v\n", err)
		return
	}

	// Configure QUIC transport
	quicConfig := &quic.Config{}

	// Start QUIC listener
	listener, err := quic.ListenAddr(":8082", tlsConfig, quicConfig)
	if err != nil {
		fmt.Printf("Failed to start QUIC server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("QUIC Server started on :8082")

	// Accept and handle connections
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		fmt.Printf("New connection from: %s\n", conn.RemoteAddr())

		// Handle each connection in a goroutine
		go func(conn quic.Connection) {
			// Accept streams
			for {
				stream, err := conn.AcceptStream(context.Background())
				if err != nil {
					fmt.Printf("Failed to accept stream: %v\n", err)
					return
				}

				// Handle each stream in a goroutine
				go handleStream(stream)
			}
		}(conn)
	}
}
