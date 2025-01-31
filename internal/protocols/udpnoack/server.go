package udpfast

import (
	"fmt"
	"net"
)

type Server struct {
	conn *net.UDPConn
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.conn = conn
	go s.handleConnections()
	return nil
}

func (s *Server) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *Server) handleConnections() {
	buffer := make([]byte, 8192)
	for {
		n, _, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			return // server closed
		}

		// Just count the messages, no processing needed
		_ = n
	}
}
