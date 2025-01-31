package grpc

import (
	"context"
	"fmt"
	"time"

	"protobench/internal/model"
	"protobench/internal/protocols/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.MessageServiceClient
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
	return "gRPC"
}

func (c *Client) SendMessage(msg *model.Message) error {
	if c.conn == nil {
		conn, err := grpc.Dial(":"+c.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
		c.conn = conn
		c.client = proto.NewMessageServiceClient(conn)
	}

	protoMsg := &proto.Message{
		Id:        msg.ID,
		Timestamp: timestamppb.New(msg.Timestamp),
		Content:   msg.Content,
		Number:    msg.Number,
		IsValid:   msg.IsValid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.client.SendMessage(ctx, protoMsg)
	return err
} 
