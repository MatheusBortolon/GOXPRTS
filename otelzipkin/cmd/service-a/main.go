package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goxprts/otelzipkin/internal/otel"
	"github.com/goxprts/otelzipkin/internal/validator"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type cepRequest struct {
	CEP string `json:"cep"`
}

func main() {
	ctx := context.Background()
	shutdown, err := otel.Setup(ctx, "service-a")
	if err != nil {
		log.Fatalf("otel setup failed: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("otel shutdown error: %v", err)
		}
	}()

	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://service-b:8081"
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/zipcode", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req cepRequest
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}

		if !validator.IsValidCEP(req.CEP) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}

		payload, err := json.Marshal(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("%s/weather", serviceBURL)
		proxyReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		proxyReq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(proxyReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		for k, values := range resp.Header {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	})

	handler := otelhttp.NewHandler(mux, "service-a")

	port := os.Getenv("SERVICE_A_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Service A running on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
