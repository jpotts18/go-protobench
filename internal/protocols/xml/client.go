package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"protobench/internal/model"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	port       string
	server     *Server
}

func NewClient(port string) *Client {
	return &Client{
		baseURL:    fmt.Sprintf("http://localhost:%s", port),
		httpClient: &http.Client{Timeout: time.Second},
		port:       port,
		server:     NewServer(port),
	}
}

func (c *Client) StartServer() error {
	return c.server.Start()
}

func (c *Client) StopServer() error {
	return c.server.Stop()
}

func (c *Client) Name() string {
	return "XML"
}

func (c *Client) SendMessage(msg *model.Message) error {
	data, err := xml.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/message", "application/xml", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s - %s", resp.Status, string(body))
	}

	return nil
} 
