package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const host = "localhost:42069"

func main() {
	resolver, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		log.Fatalf("Could not resolve UDP address: %s", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, resolver)
	if err != nil {
		log.Fatalf("Could not connect: %s", err.Error())
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading line: %s", err.Error())
		}
		written, err := conn.Write([]byte(line))
		if err != nil {
			log.Fatalf("Error writing: %s%s", err.Error(), written)
		}

		fmt.Printf("Message sent: %s", message)
	}
}
