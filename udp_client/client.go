package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Resolve the UDP address
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		fmt.Printf("Failed to resolve address: %v\n", err)
		return
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to UDP server at localhost:8081")
	fmt.Println("Type messages to send:")

	// Buffer for receiving data
	buffer := make([]byte, 1024)

	// Create scanner for user input
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		message := scanner.Text()

		// Send message to server
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}

		// Set a deadline for reading the response (2 seconds)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		// Read response from server
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timeout: No response from server")
			} else {
				fmt.Printf("Error reading from server: %v\n", err)
			}
		} else {
			fmt.Printf("Received: %s\n", string(buffer[:n]))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
