package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"labgcr/internal/app"
)

type fakeService struct {
	temps app.Temps
	err   error
}

func (f fakeService) GetTemps(ctx context.Context, cep string) (app.Temps, error) {
	return f.temps, f.err
}

func TestHandlerInvalidCEP(t *testing.T) {
	handler := NewHandler(fakeService{})

	request := httptest.NewRequest(http.MethodGet, "/weather?cep=123", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", response.Code)
	}

	if response.Body.String() != "invalid zipcode" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestHandlerCEPNotFound(t *testing.T) {
	handler := NewHandler(fakeService{err: app.ErrCEPNotFound})

	request := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}

	if response.Body.String() != "can not find zipcode" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestHandlerSuccess(t *testing.T) {
	handler := NewHandler(fakeService{temps: app.Temps{TempC: 10, TempF: 50, TempK: 283}})

	request := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body struct {
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if body.TempC != 10 || body.TempF != 50 || body.TempK != 283 {
		t.Fatalf("unexpected response: %+v", body)
	}
}
