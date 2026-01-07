package grpc

import (
	"context"
	"log"
	"net"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer   *grpc.Server
	orderService *orders.OrderService
}

func NewServer(orderService *orders.OrderService) *Server {
	return &Server{
		grpcServer:   grpc.NewServer(),
		orderService: orderService,
	}
}

func (s *Server) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening on %s", address)
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func (s *Server) ListOrders(ctx context.Context, req interface{}) (interface{}, error) {
	return s.orderService.ListOrders(ctx)
}
