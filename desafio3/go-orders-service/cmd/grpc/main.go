package main

import (
	"log"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/transport/grpc"
)

func main() {
	repo := orders.NewPostgresOrderRepository()
	ordersService := orders.NewOrderService(repo)

	server := grpc.NewServer(ordersService)
	if err := server.Start(":50051"); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
