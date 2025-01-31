package bson

import (
	"fmt"
	"net"
	"time"

	"protobench/internal/model"

	"go.mongodb.org/mongo-driver/bson"
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
	return "BSON"
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

	// Direct BSON serialization
	data, err := bson.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	c.conn.SetReadDeadline(time.Now().Add(time.Second))
	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Wait for acknowledgment
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to receive acknowledgment: %w", err)
	}

	var response struct{ Success bool }
	if err := bson.Unmarshal(buffer[:n], &response); err != nil {
		return fmt.Errorf("failed to unmarshal acknowledgment: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("server reported failure")
	}

	return nil
}
