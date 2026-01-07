package rest

import (
	"encoding/json"
	"net/http"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/orders"
	"github.com/gorilla/mux"
)

type OrderHandler struct {
	service *orders.OrderService
}

func NewOrderHandler(service *orders.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	orderList, err := h.service.ListOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderList)
}

func (h *OrderHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order orders.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderID, err := h.service.CreateOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": orderID})
}

func ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]orders.Order{})
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, "Not found", http.StatusNotFound)
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": "1"})
}

func RegisterRoutes(router *mux.Router, handler *OrderHandler) {
	router.HandleFunc("/orders", handler.ListOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/{id}", handler.GetOrderHandler).Methods("GET")
	router.HandleFunc("/orders", handler.CreateOrderHandler).Methods("POST")
}
