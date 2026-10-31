package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type serverState int

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	portString := strconv.Itoa(port)

	listener, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, fmt.Errorf("problem creating listener: %s", err)
	}
	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("unable to close listener: %s", err)
	}
	s.closed.Store(true)

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil && !s.closed.Load() {
			log.Fatalf("error accepting connection: %s", err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\r\n")
	err := conn.Close()
	if err != nil {
		log.Fatalf("error closing connection: %s", err)
	}
}
