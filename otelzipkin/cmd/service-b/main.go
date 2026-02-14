package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goxprts/otelzipkin/internal/clients"
	otelsetup "github.com/goxprts/otelzipkin/internal/otel"
	"github.com/goxprts/otelzipkin/internal/validator"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type cepRequest struct {
	CEP string `json:"cep"`
}

type weatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	ctx := context.Background()
	shutdown, err := otelsetup.Setup(ctx, "service-b")
	if err != nil {
		log.Fatalf("otel setup failed: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("otel shutdown error: %v", err)
		}
	}()

	weatherKey := os.Getenv("WEATHER_API_KEY")
	if weatherKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	viaCEP := clients.NewViaCEPClient(os.Getenv("VIACEP_BASE_URL"))
	weather := clients.NewWeatherAPIClient(weatherKey, os.Getenv("WEATHER_API_BASE_URL"))

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req cepRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}

		if !validator.IsValidCEP(req.CEP) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}

		tracer := otel.Tracer("service-b")
		ctx, span := tracer.Start(r.Context(), "lookup_zipcode")
		city, err := viaCEP.LookupCity(ctx, req.CEP)
		span.End()
		if err != nil {
			if errors.Is(err, clients.ErrZipcodeNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("can not find zipcode"))
				return
			}
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		ctx, span = tracer.Start(r.Context(), "lookup_weather")
		tempC, err := weather.GetTemperatureC(ctx, city)
		span.End()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		resp := weatherResponse{
			City:  city,
			TempC: tempC,
			TempF: tempC*1.8 + 32,
			TempK: tempC + 273,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	handler := otelhttp.NewHandler(mux, "service-b")

	port := os.Getenv("SERVICE_B_PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Service B running on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
