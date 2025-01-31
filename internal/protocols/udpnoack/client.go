package udpfast

import (
	"fmt"
	"net"

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
	return "UDP-Fast"
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

	// Format message as ASCII
	msgBody := MessageBody{
		ID:      msg.ID,
		Content: msg.Content,
		Number:  msg.Number,
	}
	asciiMsg := FormatMessage(msgBody)

	// Add protocol header with checksum
	data, err := EncodeMessage([]byte(asciiMsg))
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	// Fire and forget - no acknowledgment
	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
} 
