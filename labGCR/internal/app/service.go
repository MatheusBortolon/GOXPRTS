package app

import (
	"context"
	"errors"
)

type Temps struct {
	TempC float64
	TempF float64
	TempK float64
}

type CEPClient interface {
	LookupCity(ctx context.Context, cep string) (string, error)
}

type WeatherClient interface {
	CurrentTempC(ctx context.Context, city string) (float64, error)
}

type Service interface {
	GetTemps(ctx context.Context, cep string) (Temps, error)
}

var (
	ErrCEPNotFound   = errors.New("can not find zipcode")
	ErrMissingAPIKey = errors.New("missing weather api key")
)

type service struct {
	cepClient     CEPClient
	weatherClient WeatherClient
}

func NewService(cepClient CEPClient, weatherClient WeatherClient) Service {
	return &service{cepClient: cepClient, weatherClient: weatherClient}
}

func (s *service) GetTemps(ctx context.Context, cep string) (Temps, error) {
	city, err := s.cepClient.LookupCity(ctx, cep)
	if err != nil {
		return Temps{}, err
	}

	tempC, err := s.weatherClient.CurrentTempC(ctx, city)
	if err != nil {
		return Temps{}, err
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return Temps{TempC: tempC, TempF: tempF, TempK: tempK}, nil
}
