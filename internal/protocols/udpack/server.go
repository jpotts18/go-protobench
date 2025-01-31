package udp

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Server struct {
	conn     *net.UDPConn
	port     string
	messages map[uint64]*messageAssembler
}

type messageAssembler struct {
	chunks    map[uint32][]byte
	total     uint32
	completed bool
}

func NewServer(port string) *Server {
	return &Server{
		port:     port,
		messages: make(map[uint64]*messageAssembler),
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
	buffer := make([]byte, maxChunkSize+headerSize)
	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			return
		}

		if n < headerSize {
			continue
		}

		// Extract header info
		seqNum := binary.BigEndian.Uint64(buffer[0:8])
		chunkNum := binary.BigEndian.Uint32(buffer[8:12])
		totalChunks := binary.BigEndian.Uint32(buffer[12:16])

		// Send acknowledgment
		s.conn.WriteToUDP(buffer[:headerSize], remoteAddr)

		// Store chunk
		_, exists := s.messages[seqNum]
		if !exists {
			s.messages[seqNum] = &messageAssembler{
				chunks: make(map[uint32][]byte),
				total:  totalChunks,
			}
		}
		s.messages[seqNum].chunks[chunkNum] = append([]byte{}, buffer[headerSize:n]...)
	}
}
