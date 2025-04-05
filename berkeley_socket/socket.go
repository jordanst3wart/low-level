package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	// Create a socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Printf("Socket creation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Client socket created with file descriptor: %d\n", fd)

	// Define the server address to connect to
	serverAddr := syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1}, // 127.0.0.1
	}

	// Connect to the server (equivalent to connect() call)
	err = syscall.Connect(fd, &serverAddr)
	if err != nil {
		fmt.Printf("Connect failed: %v\n", err)
		syscall.Close(fd)
		os.Exit(1)
	}
	fmt.Println("Connected to server at 127.0.0.1:8080")

	// Example: send a message to the server
	message := []byte("Hello from client!")
	// Use syscall.Write instead of syscall.Send
	_, err = syscall.Write(fd, message)
	if err != nil {
		fmt.Printf("Send failed: %v\n", err)
		syscall.Close(fd)
		os.Exit(1)
	}
	fmt.Println("Message sent to server")

	// Example: receive a response
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		fmt.Printf("Read failed: %v\n", err)
	} else {
		fmt.Printf("Received %d bytes from server: %s\n", n, string(buf[:n]))
	}

	// Close the socket
	syscall.Close(fd)
	fmt.Println("Socket closed")
}
