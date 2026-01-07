package orders

import (
	"context"
	"testing"
)

func TestNewOrderService(t *testing.T) {
	repo := &MockOrderRepository{}
	service := NewOrderService(repo)

	if service == nil {
		t.Error("NewOrderService returned nil")
	}

	if service.repo != repo {
		t.Error("OrderService repo not set correctly")
	}
}

func TestCreateOrderWithValidAmount(t *testing.T) {
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			return "123", nil
		},
	}

	service := NewOrderService(repo)
	order := &Order{
		ID:     "123",
		Amount: 50.0,
		Status: "pending",
	}

	id, err := service.CreateOrder(order)

	if err != nil {
		t.Errorf("CreateOrder returned error: %v", err)
	}

	if id != "123" {
		t.Errorf("Expected id '123', got '%s'", id)
	}
}

func TestCreateOrderWithInvalidAmount(t *testing.T) {
	repo := &MockOrderRepository{}
	service := NewOrderService(repo)

	tests := []float64{0, -10.0, -0.01}

	for _, amount := range tests {
		order := &Order{
			ID:     "123",
			Amount: amount,
			Status: "pending",
		}

		_, err := service.CreateOrder(order)

		if err == nil {
			t.Errorf("CreateOrder with amount %v should have returned error", amount)
		}
	}
}

func TestGetOrder(t *testing.T) {
	expectedOrder := &Order{
		ID:         "123",
		Amount:     50.0,
		Status:     "pending",
		CustomerID: "cust_1",
	}

	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			if id == "123" {
				return expectedOrder, nil
			}
			return nil, nil
		},
	}

	service := NewOrderService(repo)
	order, err := service.GetOrder(context.Background(), "123")

	if err != nil {
		t.Errorf("GetOrder returned error: %v", err)
	}

	if order == nil {
		t.Error("GetOrder returned nil order")
	} else if order.ID != "123" {
		t.Errorf("Expected order id '123', got '%s'", order.ID)
	}
}

func TestListOrders(t *testing.T) {
	expectedOrders := []*Order{
		{ID: "1", Amount: 10.0, Status: "pending"},
		{ID: "2", Amount: 20.0, Status: "completed"},
	}

	repo := &MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*Order, error) {
			return expectedOrders, nil
		},
	}

	service := NewOrderService(repo)
	orders, err := service.ListOrders(context.Background())

	if err != nil {
		t.Errorf("ListOrders returned error: %v", err)
	}

	if len(orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(orders))
	}
}

func TestListOrdersEmpty(t *testing.T) {
	repo := &MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*Order, error) {
			return []*Order{}, nil
		},
	}

	service := NewOrderService(repo)
	orders, err := service.ListOrders(context.Background())

	if err != nil {
		t.Errorf("ListOrders returned error: %v", err)
	}

	if len(orders) != 0 {
		t.Errorf("Expected 0 orders, got %d", len(orders))
	}
}
