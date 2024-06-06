package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Connecting to server...")

	// Connect to server
	//conn, err := net.Dial("tcp", "10.177.240.82:8080")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Create a new file to save received data
	file, err := os.Create("aliaa.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Receive and write file data
	buffer := make([]byte, 1024)
	for {
		// Read from the connection
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}

		// If EOF reached, break the loop
		if n == 0 {
			break
		}

		// Write to the file
		_, err = file.Write(buffer[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("File received!")
}
