package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/gorilla/mux"
)

func TestListOrdersHandler(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return []*orders.Order{
				{ID: "1", Amount: 10.0, Status: "pending"},
				{ID: "2", Amount: 20.0, Status: "completed"},
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

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}
}

func TestListOrdersHandlerError(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		ListFunc: func(ctx context.Context) ([]*orders.Order, error) {
			return nil, nil
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
}

func TestGetOrderHandler(t *testing.T) {
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
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("GET", "/orders/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetOrderHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetOrderHandlerNotFound(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*orders.Order, error) {
			return nil, nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("GET", "/orders/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetOrderHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestCreateOrderHandler(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{
		CreateFunc: func(order *orders.Order) (string, error) {
			return "new-id", nil
		},
	}

	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	order := orders.Order{
		ID:         "1",
		Amount:     99.99,
		Status:     "pending",
		CustomerID: "cust_123",
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

func TestCreateOrderHandlerBadRequest(t *testing.T) {
	mockRepo := &orders.MockOrderRepository{}
	service := orders.NewOrderService(mockRepo)
	handler := NewOrderHandler(service)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateOrderHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
