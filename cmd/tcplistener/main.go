package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	fmt.Println("Setting up listener")
	fmt.Println("=====================================")

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not open network: %s\n", err)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("problem with connection: %s\n", err)
		}

		fmt.Printf("Connection has been accepted at %s\n", connection.RemoteAddr())
		fmt.Println("=====================================")

		lines := getLinesChannel(connection)
		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	line := make(chan string)
	go func() {
		currentLine := ""
		defer f.Close()
		defer close(line)
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					line <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			section := string(buffer[:n])
			parts := strings.Split(section, "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				line <- currentLine
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return line
}
