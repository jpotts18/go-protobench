package udp

import (
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

	msgBody := MessageBody{
		ID:      msg.ID,
		Content: msg.Content,
		Number:  msg.Number,
	}
	asciiMsg := FormatMessage(msgBody)
	data, err := EncodeMessage([]byte(asciiMsg))
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	// Set a very short timeout
	c.conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))

	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Wait for acknowledgment with timeout
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// On timeout, just return without error
			return nil
		}
		return fmt.Errorf("failed to receive acknowledgment: %w", err)
	}

	// Quick verify of acknowledgment without decoding
	if n < 2 || string(buffer[:2]) != "ok" {
		return fmt.Errorf("invalid acknowledgment")
	}

	return nil
}
