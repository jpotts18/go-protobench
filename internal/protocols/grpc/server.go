package grpc

import (
	"context"
	"fmt"
	"net"

	"protobench/internal/protocols/grpc/proto"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	port   string
	proto.UnimplementedMessageServiceServer
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.server = grpc.NewServer()
	proto.RegisterMessageServiceServer(s.server, s)
	
	go s.server.Serve(lis)
	return nil
}

func (s *Server) Stop() error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

func (s *Server) SendMessage(ctx context.Context, msg *proto.Message) (*proto.Response, error) {
	return &proto.Response{
		Success: true,
		Message: "Message received",
	}, nil
} 
