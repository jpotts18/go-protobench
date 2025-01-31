package bson

import (
	"encoding/binary"
	"fmt"
	"net"

	"protobench/internal/model"

	"go.mongodb.org/mongo-driver/bson"
)

type Client struct {
	conn   net.Conn
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
		conn, err := net.Dial("tcp", ":"+c.port)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		c.conn = conn
	}

	data, err := bson.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send length prefix
	size := uint32(len(data))
	if err := binary.Write(c.conn, binary.BigEndian, size); err != nil {
		return fmt.Errorf("failed to send size: %w", err)
	}

	// Send data
	if _, err := c.conn.Write(data); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
