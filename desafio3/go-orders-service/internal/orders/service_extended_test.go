package orders

import (
	"context"
	"errors"
	"testing"
)

func TestNewOrderServiceInit(t *testing.T) {
	repo := &MockOrderRepository{}
	service := NewOrderService(repo)

	if service == nil {
		t.Error("NewOrderService returned nil")
	}

	if service.repo != repo {
		t.Error("OrderService repo not set correctly")
	}
}

func TestCreateOrderValidAmountSuccess(t *testing.T) {
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

func TestCreateOrderAmountRangeValidation(t *testing.T) {
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

func TestCreateOrderRepositoryError(t *testing.T) {
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			return "", errors.New("database error")
		},
	}

	service := NewOrderService(repo)
	order := &Order{
		ID:     "123",
		Amount: 50.0,
		Status: "pending",
	}

	_, err := service.CreateOrder(order)

	if err == nil {
		t.Error("CreateOrder should have returned error from repository")
	}
}

func TestGetOrderSuccess(t *testing.T) {
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

func TestGetOrderNotFoundCase(t *testing.T) {
	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			return nil, nil
		},
	}

	service := NewOrderService(repo)
	order, err := service.GetOrder(context.Background(), "999")

	if err != nil {
		t.Errorf("GetOrder returned error: %v", err)
	}

	if order != nil {
		t.Error("GetOrder should return nil for non-existent order")
	}
}

func TestGetOrderErrorHandling(t *testing.T) {
	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewOrderService(repo)
	order, err := service.GetOrder(context.Background(), "123")

	if err == nil {
		t.Error("GetOrder should have returned error")
	}

	if order != nil {
		t.Error("GetOrder should return nil on error")
	}
}

func TestListOrdersMultiple(t *testing.T) {
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

func TestListOrdersEmptyResult(t *testing.T) {
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

func TestListOrdersRepositoryError(t *testing.T) {
	repo := &MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*Order, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewOrderService(repo)
	orders, err := service.ListOrders(context.Background())

	if err == nil {
		t.Error("ListOrders should have returned error")
	}

	if orders != nil {
		t.Error("ListOrders should return nil on error")
	}
}

func TestListOrdersSingleItem(t *testing.T) {
	expectedOrders := []*Order{
		{ID: "1", Amount: 100.0, Status: "completed", CustomerID: "cust_1"},
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

	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}

	if orders[0].ID != "1" {
		t.Errorf("Expected order id '1', got '%s'", orders[0].ID)
	}
}

func TestCreateOrderLargeAmount(t *testing.T) {
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			return "large_order", nil
		},
	}

	service := NewOrderService(repo)
	order := &Order{
		ID:     "large_order",
		Amount: 999999.99,
		Status: "pending",
	}

	id, err := service.CreateOrder(order)

	if err != nil {
		t.Errorf("CreateOrder with large amount returned error: %v", err)
	}

	if id != "large_order" {
		t.Errorf("Expected id 'large_order', got '%s'", id)
	}
}
func TestCreateOrderMinimumValidAmount(t *testing.T) {
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			return "minimum-id", nil
		},
	}

	service := NewOrderService(repo)
	order := &Order{
		ID:     "min",
		Amount: 0.01,
		Status: "pending",
	}

	id, err := service.CreateOrder(order)

	if err != nil {
		t.Errorf("CreateOrder with minimum valid amount should succeed: %v", err)
	}

	if id != "minimum-id" {
		t.Errorf("Expected id 'minimum-id', got '%s'", id)
	}
}

func TestGetOrderWithContext(t *testing.T) {
	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return &Order{ID: id, Amount: 75.0, Status: "pending"}, nil
		},
	}

	service := NewOrderService(repo)
	ctx := context.Background()

	order, err := service.GetOrder(ctx, "123")

	if err != nil {
		t.Errorf("GetOrder returned error: %v", err)
	}

	if order == nil || order.ID != "123" {
		t.Error("GetOrder should return order with matching ID")
	}
}

func TestListOrdersWithContext(t *testing.T) {
	repo := &MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*Order, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return []*Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
				{ID: "2", Amount: 20.0, Status: "completed"},
				{ID: "3", Amount: 30.0, Status: "pending"},
			}, nil
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
}

func TestCreateOrderRepositoryCallValidation(t *testing.T) {
	receivedOrder := (*Order)(nil)
	repo := &MockOrderRepository{
		CreateFunc: func(order *Order) (string, error) {
			receivedOrder = order
			return "validated-id", nil
		},
	}

	service := NewOrderService(repo)
	testOrder := &Order{
		ID:         "test-id",
		CustomerID: "cust-123",
		Amount:     55.5,
		Status:     "pending",
	}

	id, err := service.CreateOrder(testOrder)

	if err != nil {
		t.Errorf("CreateOrder should not return error: %v", err)
	}

	if id != "validated-id" {
		t.Errorf("Expected id 'validated-id', got '%s'", id)
	}

	if receivedOrder == nil {
		t.Error("Repository.Create should have been called with order")
	}
}

func TestListOrdersLargeDataSet(t *testing.T) {
	var orderList []*Order
	for i := 0; i < 50; i++ {
		orderList = append(orderList, &Order{
			ID:     "order-" + string(rune(48+i%10)),
			Amount: float64(i) * 10.5,
			Status: "pending",
		})
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

	if len(orders) != 50 {
		t.Errorf("Expected 50 orders, got %d", len(orders))
	}
}

func TestGetOrderNotFound(t *testing.T) {
	repo := &MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*Order, error) {
			return nil, nil
		},
	}

	service := NewOrderService(repo)
	ctx := context.Background()

	order, err := service.GetOrder(ctx, "nonexistent")

	if err != nil {
		t.Errorf("GetOrder should not return error for nonexistent: %v", err)
	}

	if order != nil {
		t.Error("GetOrder should return nil for nonexistent order")
	}
}
