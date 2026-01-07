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

func TestCreateOrderHandlerInvalidAmountError(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		CreateFunc: func(order *orders.Order) (string, error) {
			return "", errors.New("invalid amount")
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	order := orders.Order{
		ID:     "1",
		Amount: 0,
		Status: "pending",
	}

	body, _ := json.Marshal(order)
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestListOrdersHandlerWithMultipleResults(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
				{ID: "2", Amount: 20.0, Status: "completed"},
				{ID: "3", Amount: 30.0, Status: "shipped"},
			}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	handler.ListOrdersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var orderList []*orders.Order
	json.NewDecoder(w.Body).Decode(&orderList)

	if len(orderList) != 3 {
		t.Errorf("Expected 3 orders in response, got %d", len(orderList))
	}
}

func TestCreateOrderHandlerSuccessfulCreation(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		CreateFunc: func(order *orders.Order) (string, error) {
			return "order_123", nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	order := orders.Order{
		ID:         "order_123",
		Amount:     99.99,
		Status:     "pending",
		CustomerID: "cust_456",
	}

	body, _ := json.Marshal(order)
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetOrderHandlerStatusOK(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			if id == "1" {
				return &orders.Order{ID: "1", Amount: 100.0, Status: "completed"}, nil
			}
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req, _ := http.NewRequest("GET", "/orders/1", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestListOrdersHandlerEmptyJSON(t *testing.T) {
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
}

func TestCreateOrderHandlerLargeAmount(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		CreateFunc: func(order *orders.Order) (string, error) {
			return "large-id", nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	order := orders.Order{CustomerID: "cust1", Amount: 999999.99, Status: "pending"}
	body, _ := json.Marshal(order)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetOrderHandlerMultipleRequests(t *testing.T) {
	callCount := 0
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			callCount++
			return &orders.Order{ID: id, Amount: 50.0, Status: "pending"}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")

	for i := 1; i <= 3; i++ {
		req, _ := http.NewRequest("GET", "/orders/id"+string(rune(48+i)), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i, w.Code)
		}
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls to GetByID, got %d", callCount)
	}
}

func TestListOrdersHandlerErrorHandling(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return nil, errors.New("database error")
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	handler.ListOrdersHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestCreateOrderHandlerMinimumAmount(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		CreateFunc: func(order *orders.Order) (string, error) {
			return "min-id", nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	order := orders.Order{CustomerID: "cust1", Amount: 0.01, Status: "pending"}
	body, _ := json.Marshal(order)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestGetOrderHandlerByID(t *testing.T) {
	expectedOrder := &orders.Order{
		ID:         "order_1",
		Amount:     50.0,
		Status:     "pending",
		CustomerID: "cust_1",
	}

	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			if id == "order_1" {
				return expectedOrder, nil
			}
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("GET", "/orders/order_1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "order_1"})
	w := httptest.NewRecorder()

	handler.GetOrderHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var order orders.Order
	json.NewDecoder(w.Body).Decode(&order)

	if order.ID != "order_1" {
		t.Errorf("Expected order id 'order_1', got '%s'", order.ID)
	}
}

func TestListOrdersHandlerResponseContentType(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	handler.ListOrdersHandler(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestCreateOrderHandlerEmptyBodyError(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetOrderHandlerResponseType(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			return &orders.Order{ID: "1", Amount: 10.0}, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("GET", "/orders/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetOrderHandler(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestCreateOrderHandlerStatusCreated(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		CreateFunc: func(order *orders.Order) (string, error) {
			return "new_id", nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	order := orders.Order{ID: "new_id", Amount: 100.0}
	body, _ := json.Marshal(order)
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}
