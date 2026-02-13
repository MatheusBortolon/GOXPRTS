package app

import (
	"context"
	"testing"
)

type fakeCEPClient struct {
	city string
	err  error
}

func (f fakeCEPClient) LookupCity(ctx context.Context, cep string) (string, error) {
	return f.city, f.err
}

type fakeWeatherClient struct {
	temp float64
	err  error
}

func (f fakeWeatherClient) CurrentTempC(ctx context.Context, city string) (float64, error) {
	return f.temp, f.err
}

func TestServiceTemps(t *testing.T) {
	service := NewService(fakeCEPClient{city: "Sao Paulo"}, fakeWeatherClient{temp: 10})
	temps, err := service.GetTemps(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if temps.TempC != 10 {
		t.Fatalf("expected tempC 10, got %v", temps.TempC)
	}
	if temps.TempF != 50 {
		t.Fatalf("expected tempF 50, got %v", temps.TempF)
	}
	if temps.TempK != 283 {
		t.Fatalf("expected tempK 283, got %v", temps.TempK)
	}
}

func TestServiceCEPNotFound(t *testing.T) {
	service := NewService(fakeCEPClient{err: ErrCEPNotFound}, fakeWeatherClient{})
	_, err := service.GetTemps(context.Background(), "01001000")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != ErrCEPNotFound {
		t.Fatalf("expected ErrCEPNotFound, got %v", err)
	}
}
