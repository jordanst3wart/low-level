package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	// Create a scanner to read from the connection
	scanner := bufio.NewScanner(conn)

	// Send welcome message
	conn.Write([]byte("Welcome to the TCP server! Type 'quit' to exit.\n"))

	// Process client messages
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Received: %s\n", message)

		// Check if client wants to quit
		if strings.ToLower(message) == "quit" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}

		// Echo the message back with a prefix
		response := fmt.Sprintf("Server: You said '%s'\n", message)
		conn.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from client: %v\n", err)
	}

	fmt.Printf("Connection closed: %s\n", conn.RemoteAddr().String())
}

func main() {
	// Listen on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP Server started on :8080")

	// Accept connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn)
	}
}
