package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrZipcodeNotFound = errors.New("zipcode not found")

type ViaCEPClient struct {
	baseURL string
	client  *http.Client
}

func NewViaCEPClient(baseURL string) *ViaCEPClient {
	if baseURL == "" {
		baseURL = "https://viacep.com.br"
	}
	return &ViaCEPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

func (c *ViaCEPClient) LookupCity(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.baseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create viacep request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("viacep request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep status: %d", resp.StatusCode)
	}

	var payload viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode viacep response: %w", err)
	}

	if payload.Erro != "" || payload.Localidade == "" {
		return "", ErrZipcodeNotFound
	}

	return payload.Localidade, nil
}
