package main

import (
	"log"
	"net/http"

	"github.com/MatheusBortolon/GOXPRTS/desafio3/go-orders-service/internal/transport/rest"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/orders", rest.ListOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/{id}", rest.GetOrderHandler).Methods("GET")
	router.HandleFunc("/orders", rest.CreateOrderHandler).Methods("POST")
	log.Println("Starting REST API server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
