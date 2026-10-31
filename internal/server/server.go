package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"httpfromtcp/internal/response"
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
			fmt.Printf("error accepting connection: %v", err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response.WriteStatusLine(conn, response.OK)
	defaultHeaders := response.GetDefaultHeaders(0)
	err := response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		fmt.Printf("error writing headers: %v", err)
	}
}
