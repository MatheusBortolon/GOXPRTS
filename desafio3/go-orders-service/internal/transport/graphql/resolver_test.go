package graphql

import (
	"context"
	"testing"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/graphql-go/graphql"
)

func TestNewResolver(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)

	resolver := NewResolver(service)

	if resolver == nil {
		t.Error("NewResolver returned nil")
	}

	if resolver.orderService != service {
		t.Error("Resolver orderService not set correctly")
	}
}

func TestListOrdersResolver(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
				{ID: "2", Amount: 20.0, Status: "completed"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	resolver := NewResolver(service)

	params := graphql.ResolveParams{
		Context: context.Background(),
	}

	result, err := resolver.ListOrdersResolver(params)

	if err != nil {
		t.Errorf("ListOrdersResolver returned error: %v", err)
	}

	if result == nil {
		t.Error("ListOrdersResolver returned nil result")
	}
}

func TestGetOrderResolver(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			if id == "1" {
				return &orders.Order{
					ID:         "1",
					Amount:     50.0,
					Status:     "pending",
					CustomerID: "cust_1",
				}, nil
			}
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	resolver := NewResolver(service)

	params := graphql.ResolveParams{
		Context: context.Background(),
		Args: map[string]interface{}{
			"id": "1",
		},
	}

	result, err := resolver.GetOrderResolver(params)

	if err != nil {
		t.Errorf("GetOrderResolver returned error: %v", err)
	}

	if result == nil {
		t.Error("GetOrderResolver returned nil result")
	}
}
