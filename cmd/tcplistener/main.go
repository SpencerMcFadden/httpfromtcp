package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/request"
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

		request, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatalf("error reading request: %v", err)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}
