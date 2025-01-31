package udp

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"protobench/internal/model"
)

type Client struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	port   string
	server *Server
}

func NewClient(port string) *Client {
	return &Client{
		port:   port,
		server: NewServer(port),
	}
}

func (c *Client) StartServer() error {
	return c.server.Start()
}

func (c *Client) StopServer() error {
	return c.server.Stop()
}

func (c *Client) Name() string {
	return "UDP"
}

const (
	maxChunkSize = 1400 // Leave room for headers
	headerSize   = 16   // 8 bytes seq + 8 bytes chunk info
)

func (c *Client) SendMessage(msg *model.Message) error {
	if c.conn == nil {
		addr, err := net.ResolveUDPAddr("udp", ":"+c.port)
		if err != nil {
			return fmt.Errorf("failed to resolve address: %w", err)
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return fmt.Errorf("failed to dial: %w", err)
		}
		c.conn = conn
		c.addr = addr
	}

	content := []byte(msg.Content)
	totalChunks := (len(content) + maxChunkSize - 1) / maxChunkSize

	// Send each chunk with retries
	for chunk := 0; chunk < totalChunks; chunk++ {
		start := chunk * maxChunkSize
		end := start + maxChunkSize
		if end > len(content) {
			end = len(content)
		}

		// Header: sequence number (8 bytes) + chunk number (4 bytes) + total chunks (4 bytes)
		header := make([]byte, headerSize)
		binary.BigEndian.PutUint64(header[0:8], uint64(msg.Number))
		binary.BigEndian.PutUint32(header[8:12], uint32(chunk))
		binary.BigEndian.PutUint32(header[12:16], uint32(totalChunks))

		// Combine header and chunk data
		data := append(header, content[start:end]...)

		// Try to send chunk with retries
		maxRetries := 3
		success := false
		for retry := 0; retry < maxRetries; retry++ {
			if err := c.sendChunkWithAck(data); err == nil {
				success = true
				break
			}
		}
		if !success {
			return fmt.Errorf("failed to send chunk %d/%d after retries", chunk+1, totalChunks)
		}
	}

	return nil
}

func (c *Client) sendChunkWithAck(data []byte) error {
	if _, err := c.conn.Write(data); err != nil {
		return err
	}

	// Wait for acknowledgment
	c.conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	ackBuf := make([]byte, headerSize)
	n, err := c.conn.Read(ackBuf)
	if err != nil || n != headerSize {
		return fmt.Errorf("ack error: %v", err)
	}

	return nil
}
