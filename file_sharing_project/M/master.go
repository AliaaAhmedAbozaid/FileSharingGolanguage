package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	chunkSize = 1024 // Adjust chunk size as needed
	timeout   = 120 * time.Second
)

func sendFileChunks(conn net.Conn, file *os.File, start, end int64) error {
	_, err := file.Seek(start, 0)
	if err != nil {
		return fmt.Errorf("error seeking file: %v", err)
	}

	buffer := make([]byte, end-start)
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	_, err = conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("error sending file chunk: %v", err)
	}
	return nil
}

func main() {
	filePath := "yarab.txt" // Path to the file to be shared

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Master node running and listening on port 8080...")

	// Accept connections until timeout
	var connections []net.Conn
	for {
		ln.(*net.TCPListener).SetDeadline(time.Now().Add(timeout)) // Set a deadline for accepting connections
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				break // Timeout reached, stop accepting connections
			}
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Slave connected:", conn.RemoteAddr())
		connections = append(connections, conn)
	}

	// Send file chunks to each client
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(len(connections))

	for i, conn := range connections {
		startOffset := int64(i) * chunkSize
		endOffset := startOffset + chunkSize

		if i == len(connections)-1 {
			endOffset = fileSize
		}

		go func(conn net.Conn, start, end int64) {
			defer conn.Close()
			err := sendFileChunks(conn, file, start, end)
			if err != nil {
				fmt.Println(err)
			}
		}(conn, startOffset, endOffset)
	}

	fmt.Println("File distribution complete.")
}
