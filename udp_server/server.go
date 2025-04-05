package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen on UDP port 8081
	addr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		fmt.Printf("Failed to resolve address: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP Server started on :8081")

	// Buffer for receiving data
	buffer := make([]byte, 1024)

	for {
		// Read message from client
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Received from %s: %s\n", clientAddr.String(), message)

		// Prepare response
		response := fmt.Sprintf("Server: You said '%s'", message)

		// Send response back to the client
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Printf("Failed to send response: %v\n", err)
		}
	}
}
