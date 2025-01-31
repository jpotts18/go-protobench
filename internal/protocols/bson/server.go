package bson

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"protobench/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	listener net.Listener
	port     string
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = listener

	go s.handleConnections()
	return nil
}

func (s *Server) handleConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Read message size
		var size uint32
		if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
			return
		}

		// Read message data
		data := make([]byte, size)
		if _, err := io.ReadFull(conn, data); err != nil {
			return
		}

		var msg model.Message
		if err := bson.Unmarshal(data, &msg); err != nil {
			continue
		}

		// Send acknowledgment
		conn.Write([]byte{1})
	}
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}
