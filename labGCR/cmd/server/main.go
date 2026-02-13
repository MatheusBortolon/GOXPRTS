package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"labgcr/internal/app"
	"labgcr/internal/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	cepBase := getEnv("CEP_API_BASE", "https://viacep.com.br")
	weatherBase := getEnv("WEATHER_API_BASE", "https://api.weatherapi.com")
	apiKey := os.Getenv("WEATHER_API_KEY")

	cepClient := app.NewViaCEPClient(cepBase, httpClient)
	weatherClient := app.NewWeatherAPIClient(weatherBase, apiKey, httpClient)
	service := app.NewService(cepClient, weatherClient)
	weatherHandler := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/weather", weatherHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
