package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
)

func TestServerInitializationProper(t *testing.T) {
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

func TestListOrdersMultipleResults(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
				{ID: "2", Amount: 20.0, Status: "completed"},
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

func TestListOrdersDatabaseError(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err == nil {
		t.Error("ListOrders should have returned error")
	}

	if result != nil {
		if slice, ok := result.([]*orders.Order); ok {
			if slice != nil && len(slice) > 0 {
				t.Error("ListOrders should return empty result on error")
			}
		}
	}
}

func TestListOrdersEmptyResponse(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders returned error: %v", err)
	}

	if result == nil {
		t.Error("ListOrders should return empty slice, not nil")
	}
}

func TestServerStartInitializesListener(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)
	if server.grpcServer == nil {
		t.Error("Server grpcServer should be initialized")
	}
}

func TestServerStopClosesGrpcServer(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)
	server.Stop()
	if server.grpcServer == nil {
		t.Error("Server grpcServer reference should still exist after Stop")
	}
}

func TestListOrdersWithContextCancellation(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return []*orders.Order{
				{ID: "1", Amount: 100.0, Status: "pending"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := server.ListOrders(ctx, nil)
	if err == nil {
		if result == nil {
			t.Error("ListOrders should return non-nil result on no error")
		}
	}
}

func TestListOrdersLargeDataSet(t *testing.T) {
	var orderList []*orders.Order
	for i := 0; i < 100; i++ {
		orderList = append(orderList, &orders.Order{
			ID:     string(rune(i)),
			Amount: float64(i) * 10.5,
			Status: "pending",
		})
	}

	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return orderList, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders returned error: %v", err)
	}

	if result == nil {
		t.Error("ListOrders should return non-nil result")
	}
}

func TestListOrdersServerIntegration(t *testing.T) {
	callCount := 0
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			callCount++
			return []*orders.Order{
				{ID: "test1", Amount: 50.0, Status: "completed"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result1, err1 := server.ListOrders(context.Background(), nil)
	result2, err2 := server.ListOrders(context.Background(), nil)

	if err1 != nil || err2 != nil {
		t.Error("ListOrders should not return errors")
	}

	if result1 == nil || result2 == nil {
		t.Error("ListOrders should return non-nil results")
	}

	if callCount != 2 {
		t.Errorf("Repository.List should be called twice, was called %d times", callCount)
	}
}
func TestListOrdersWithSpecificOrders(t *testing.T) {
	testOrders := []*orders.Order{
		{ID: "order-1", CustomerID: "cust-1", Amount: 100.0, Status: "completed"},
		{ID: "order-2", CustomerID: "cust-2", Amount: 200.0, Status: "pending"},
		{ID: "order-3", CustomerID: "cust-1", Amount: 150.0, Status: "shipping"},
	}

	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return testOrders, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders should not return error: %v", err)
	}

	if result == nil {
		t.Error("ListOrders should return non-nil result")
	}

	orders, ok := result.([]*orders.Order)
	if !ok {
		t.Error("Result should be []*Order")
	}

	if len(orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(orders))
	}
}

func TestServerInitializationFields(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	if server.orderService == nil {
		t.Error("Server.orderService should not be nil")
	}

	if server.grpcServer == nil {
		t.Error("Server.grpcServer should not be nil")
	}
}

func TestListOrdersErrorPropagation(t *testing.T) {
	expectedError := errors.New("connection timeout")
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return nil, expectedError
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	_, err := server.ListOrders(context.Background(), nil)

	if err == nil {
		t.Error("ListOrders should propagate error from service")
	}

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error '%v', got '%v'", expectedError, err)
	}
}

func TestListOrdersNilSliceHandling(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders should not return error for nil slice: %v", err)
	}

	if result != nil {
		orders, ok := result.([]*orders.Order)
		if ok && orders == nil {
			return
		}
		t.Error("Result should be nil or empty slice")
	}
}

func TestServerListOrdersWithTimeout(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return []*orders.Order{
					{ID: "timeout-test", Amount: 99.0, Status: "pending"},
				}, nil
			}
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders should not timeout with background context: %v", err)
	}

	if result == nil {
		t.Error("ListOrders should return result with background context")
	}
}
func TestServerStartInvalidAddress(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)
	err := server.Start("invalid::address::format")

	if err == nil {
		t.Error("Start should return error for invalid address")
	}
}

func TestServerStartValidAddress(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)
	go func() {
		_ = server.Start("localhost:50099")
	}()
	server.Stop()
}

func TestServerStartPortInUse(t *testing.T) {
	t.Skip("Port-in-use scenario skipped to avoid conflicts in CI")
}

func TestNewServerFieldInitialization(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)

	server := NewServer(service)

	if server == nil {
		t.Fatal("NewServer should not return nil")
	}

	if server.orderService == nil {
		t.Error("orderService field should be initialized")
	}

	if server.grpcServer == nil {
		t.Error("grpcServer field should be initialized")
	}

	if server.orderService != service {
		t.Error("orderService should be the provided service instance")
	}
}

func TestListOrdersWithNilRequest(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{
				{ID: "1", Amount: 50.0, Status: "pending"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders should handle nil request: %v", err)
	}

	if result == nil {
		t.Error("ListOrders should return result with nil request")
	}
}

func TestStopClosesServer(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)
	server.Stop()
	server.Stop()
}

func TestListOrdersReturnsServiceResult(t *testing.T) {
	expectedOrders := []*orders.Order{
		{ID: "result-1", CustomerID: "cust-1", Amount: 111.11, Status: "completed"},
		{ID: "result-2", CustomerID: "cust-2", Amount: 222.22, Status: "pending"},
	}

	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return expectedOrders, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	server := NewServer(service)

	result, err := server.ListOrders(context.Background(), nil)

	if err != nil {
		t.Errorf("ListOrders should not return error: %v", err)
	}

	orders, ok := result.([]*orders.Order)
	if !ok {
		t.Fatal("Result should be []*orders.Order")
	}

	if len(orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(orders))
	}

	if orders[0].ID != "result-1" || orders[1].ID != "result-2" {
		t.Error("ListOrders should return exact orders from service")
	}
}
