package grpc

import (
	"context"
	"testing"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
)

func TestNewServer(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)

	server := NewServer(service)

	if server == nil {
		t.Error("NewServer returned nil")
	}

	if server.orderService != service {
		t.Error("Server orderService not set correctly")
	}

	if server.grpcServer == nil {
		t.Error("Server grpcServer not initialized")
	}
}

func TestListOrders(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders returned error: %v", err)
	}

	if result == nil {
		t.Error("ListOrders returned nil result")
	}
}
