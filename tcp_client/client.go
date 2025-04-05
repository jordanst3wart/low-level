package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server at localhost:8080")

	// Create scanner for server responses
	serverScanner := bufio.NewScanner(conn)

	// Create scanner for user input
	inputScanner := bufio.NewScanner(os.Stdin)

	// Start a goroutine to read server responses
	go func() {
		for serverScanner.Scan() {
			fmt.Println(serverScanner.Text())
		}
		if err := serverScanner.Err(); err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
		}
		fmt.Println("Server connection closed")
		os.Exit(0)
	}()

	// Main loop for user input
	fmt.Println("Type messages to send (type 'quit' to exit):")
	for inputScanner.Scan() {
		message := inputScanner.Text()

		// Send message to server
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}

		// Check if user wants to quit
		if strings.ToLower(message) == "quit" {
			break
		}
	}

	if err := inputScanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
