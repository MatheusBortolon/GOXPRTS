package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/gorilla/mux"
)

func TestListOrdersHandlerWithLargeDataSet(t *testing.T) {
	var orderList []*orders.Order
	for i := 0; i < 20; i++ {
		orderList = append(orderList, &orders.Order{
			ID:     "order-" + string(rune(48+i%10)),
			Amount: float64(i) * 5.0,
			Status: "pending",
		})
	}

	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return orderList, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	handler.ListOrdersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result []orders.Order
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(result) != 20 {
		t.Errorf("Expected 20 orders in response, got %d", len(result))
	}
}

func TestGetOrderHandlerInvalidID(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")

	req, _ := http.NewRequest("GET", "/orders/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for invalid ID, got %d", w.Code)
	}
}

func TestCreateOrderHandlerDecodeError(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestListOrdersHandlerMultipleCalls(t *testing.T) {
	callCount := 0
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			callCount++
			return []*orders.Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/orders", nil)
		w := httptest.NewRecorder()

		handler.ListOrdersHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Iteration %d: Expected status 200, got %d", i, w.Code)
		}
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls to ListOrders, got %d", callCount)
	}
}

func TestGetOrderHandlerContextPropagation(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return &orders.Order{ID: id, Amount: 50.0, Status: "pending"}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")

	req, _ := http.NewRequest("GET", "/orders/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateOrderHandlerWithVariousAmounts(t *testing.T) {
	testAmounts := []float64{0.01, 10.0, 100.0, 999999.99}

	for _, amount := range testAmounts {
		mockRepo := &orders.MockOrderRepository{
			CreateFunc: func(order *orders.Order) (string, error) {
				return "created-id", nil
			},
		}

		service := orders.NewOrderService(mockRepo)
		handler := NewOrderHandler(service)

		order := orders.Order{CustomerID: "cust", Amount: amount, Status: "pending"}
		body, _ := json.Marshal(order)

		req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateOrderHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Amount %f: Expected status 201, got %d", amount, w.Code)
		}
	}
}

func TestGetOrderHandlerResponseBody(t *testing.T) {
	testOrder := &orders.Order{
		ID:         "test-123",
		CustomerID: "customer-456",
		Amount:     99.99,
		Status:     "completed",
	}

	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			if id == "test-123" {
				return testOrder, nil
			}
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")

	req, _ := http.NewRequest("GET", "/orders/test-123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result orders.Order
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if result.ID != testOrder.ID {
		t.Errorf("Expected ID %s, got %s", testOrder.ID, result.ID)
	}

	if result.Amount != testOrder.Amount {
		t.Errorf("Expected amount %f, got %f", testOrder.Amount, result.Amount)
	}
}

func TestListOrdersHandlerContentTypeHeader(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	handler.ListOrdersHandler(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}

func TestStandaloneListOrdersHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	ListOrdersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var result []orders.Order
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty array, got %d items", len(result))
	}
}

func TestStandaloneGetOrderHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/orders/1", nil)
	w := httptest.NewRecorder()

	GetOrderHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestStandaloneCreateOrderHandler(t *testing.T) {
	order := orders.Order{CustomerID: "cust1", Amount: 100.0, Status: "pending"}
	body, _ := json.Marshal(order)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateOrderHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var result map[string]string
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if result["id"] != "1" {
		t.Errorf("Expected id '1', got '%s'", result["id"])
	}
}

func TestRegisterRoutes(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{}, nil
		},
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			return &orders.Order{ID: id, Amount: 100.0, Status: "pending"}, nil
		},
		CreateFunc: func(order *orders.Order) (string, error) {
			return "created-id", nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	router := mux.NewRouter()
	RegisterRoutes(router, handler)

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ListOrders route: Expected status 200, got %d", w.Code)
	}

	req2, _ := http.NewRequest("GET", "/orders/123", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("GetOrder route: Expected status 200, got %d", w2.Code)
	}

	order := orders.Order{CustomerID: "cust", Amount: 100.0, Status: "pending"}
	body, _ := json.Marshal(order)
	req3, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusCreated {
		t.Errorf("CreateOrder route: Expected status 201, got %d", w3.Code)
	}
}

func TestGetOrderHandlerServiceError(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			return nil, errors.New("database error")
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")

	req, _ := http.NewRequest("GET", "/orders/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 on service error, got %d", w.Code)
	}
}
