package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ViaCEPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewViaCEPClient(baseURL string, httpClient *http.Client) *ViaCEPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &ViaCEPClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *ViaCEPClient) LookupCity(ctx context.Context, cep string) (string, error) {
	endpoint := fmt.Sprintf("%s/ws/%s/json/", c.baseURL, cep)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep status: %d", response.StatusCode)
	}

	var payload struct {
		Localidade string `json:"localidade"`
		Erro       bool   `json:"erro"`
	}

	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}

	if payload.Erro || payload.Localidade == "" {
		return "", ErrCEPNotFound
	}

	return payload.Localidade, nil
}

type WeatherAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewWeatherAPIClient(baseURL, apiKey string, httpClient *http.Client) *WeatherAPIClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &WeatherAPIClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *WeatherAPIClient) CurrentTempC(ctx context.Context, city string) (float64, error) {
	if c.apiKey == "" {
		return 0, ErrMissingAPIKey
	}

	endpoint := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s&aqi=no", c.baseURL, url.QueryEscape(c.apiKey), url.QueryEscape(city))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi status: %d", response.StatusCode)
	}

	var payload struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return 0, err
	}

	return payload.Current.TempC, nil
}
