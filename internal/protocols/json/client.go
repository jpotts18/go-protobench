package json

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		baseURL: fmt.Sprintf("http://localhost:%s", port),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
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
	return "JSON"
}

func (c *Client) SendMessage(msg *model.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/message",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
