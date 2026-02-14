package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WeatherAPIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewWeatherAPIClient(apiKey, baseURL string) *WeatherAPIClient {
	if baseURL == "" {
		baseURL = "https://api.weatherapi.com"
	}
	return &WeatherAPIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *WeatherAPIClient) GetTemperatureC(ctx context.Context, city string) (float64, error) {
	endpoint := fmt.Sprintf("%s/v1/current.json", c.baseURL)
	params := url.Values{}
	params.Set("key", c.apiKey)
	params.Set("q", city)
	params.Set("aqi", "no")
	urlWithParams := endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithParams, nil)
	if err != nil {
		return 0, fmt.Errorf("create weatherapi request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("weatherapi request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi status: %d", resp.StatusCode)
	}

	var payload weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("decode weatherapi response: %w", err)
	}

	return payload.Current.TempC, nil
}
