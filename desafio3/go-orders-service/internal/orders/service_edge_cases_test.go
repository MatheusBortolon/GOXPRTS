package orders

import (
	"context"
	"errors"
	"testing"
)

func TestCreateOrderZeroAmount(t *testing.T) {
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			return "", errors.New("invalid amount")
		},
	}

	service := NewOrderService(repo)
	order := &Order{
		ID:     "zero",
		Amount: 0,
		Status: "pending",
	}

	id, err := service.CreateOrder(order)

	if err == nil {
		t.Error("CreateOrder with zero amount should fail")
	}

	if id != "" {
		t.Errorf("Expected empty id on error, got '%s'", id)
	}
}

func TestCreateOrderNegativeAmount(t *testing.T) {
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			return "", errors.New("invalid amount")
		},
	}

	service := NewOrderService(repo)
	order := &Order{
		ID:     "negative",
		Amount: -10.0,
		Status: "pending",
	}

	id, err := service.CreateOrder(order)

	if err == nil {
		t.Error("CreateOrder with negative amount should fail")
	}

	if id != "" {
		t.Errorf("Expected empty id on error, got '%s'", id)
	}
}

func TestListOrdersSortOrder(t *testing.T) {
	orderList := []*Order{
		{ID: "3", Amount: 300.0, Status: "pending"},
		{ID: "1", Amount: 100.0, Status: "completed"},
		{ID: "2", Amount: 200.0, Status: "pending"},
	}

	repo := &MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*Order, error) {
			return orderList, nil
		},
	}

	service := NewOrderService(repo)
	ctx := context.Background()

	orders, err := service.ListOrders(ctx)

	if err != nil {
		t.Errorf("ListOrders returned error: %v", err)
	}

	if len(orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(orders))
	}

	if orders[0].ID != "3" || orders[1].ID != "1" || orders[2].ID != "2" {
		t.Error("ListOrders should preserve order from repository")
	}
}

func TestGetOrderContextDeadline(t *testing.T) {
	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			return &Order{ID: id, Amount: 100.0, Status: "pending"}, nil
		},
	}

	service := NewOrderService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	order, err := service.GetOrder(ctx, "123")

	if err == nil {
		t.Error("GetOrder should return error for cancelled context")
	}

	if order != nil {
		t.Error("GetOrder should return nil order for cancelled context")
	}
}

func TestListOrdersRepositoryCallCount(t *testing.T) {
	callCount := 0
	repo := &MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*Order, error) {
			callCount++
			return []*Order{
				{ID: "1", Amount: 50.0, Status: "pending"},
			}, nil
		},
	}

	service := NewOrderService(repo)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_, err := service.ListOrders(ctx)
		if err != nil {
			t.Errorf("ListOrders iteration %d failed: %v", i, err)
		}
	}

	if callCount != 5 {
		t.Errorf("Expected 5 calls to List, got %d", callCount)
	}
}

func TestCreateOrderMultipleSuccessfulCalls(t *testing.T) {
	callCount := 0
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			callCount++
			return "success-" + order.ID, nil
		},
	}

	service := NewOrderService(repo)

	for i := 0; i < 3; i++ {
		order := &Order{
			ID:     "test-" + string(rune(48+i)),
			Amount: float64(i+1) * 10,
			Status: "pending",
		}

		id, err := service.CreateOrder(order)
		if err != nil {
			t.Errorf("CreateOrder iteration %d failed: %v", i, err)
		}

		if id != "success-test-"+string(rune(48+i)) {
			t.Errorf("Iteration %d: expected 'success-test-%s', got '%s'", i, string(rune(48+i)), id)
		}
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls to Create, got %d", callCount)
	}
}

func TestGetOrderWithDifferentIDs(t *testing.T) {
	testIDs := []string{"order-1", "order-2", "order-3"}
	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			for _, testID := range testIDs {
				if id == testID {
					return &Order{ID: id, Amount: 100.0, Status: "pending"}, nil
				}
			}
			return nil, nil
		},
	}

	service := NewOrderService(repo)
	ctx := context.Background()

	for _, id := range testIDs {
		order, err := service.GetOrder(ctx, id)

		if err != nil {
			t.Errorf("GetOrder failed for id %s: %v", id, err)
		}

		if order == nil {
			t.Errorf("GetOrder should return order for id %s", id)
		} else if order.ID != id {
			t.Errorf("GetOrder returned order with wrong ID: expected %s, got %s", id, order.ID)
		}
	}
	unknownOrder, err := service.GetOrder(ctx, "unknown")
	if err != nil {
		t.Errorf("GetOrder for unknown ID should not return error: %v", err)
	}
	if unknownOrder != nil {
		t.Error("GetOrder should return nil for unknown id")
	}
}
