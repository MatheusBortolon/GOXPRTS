package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"labgcr/internal/app"
)

type Handler struct {
	service app.Service
}

func NewHandler(service app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cep := request.URL.Query().Get("cep")
	if !isValidCEP(cep) {
		writePlain(writer, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	temps, err := h.service.GetTemps(request.Context(), cep)
	if err != nil {
		switch {
		case errors.Is(err, app.ErrCEPNotFound):
			writePlain(writer, http.StatusNotFound, "can not find zipcode")
		case errors.Is(err, app.ErrMissingAPIKey):
			writePlain(writer, http.StatusInternalServerError, "internal server error")
		default:
			writePlain(writer, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response := struct {
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}{
		TempC: temps.TempC,
		TempF: temps.TempF,
		TempK: temps.TempK,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func isValidCEP(value string) bool {
	if len(value) != 8 {
		return false
	}

	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func writePlain(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(message))
}
