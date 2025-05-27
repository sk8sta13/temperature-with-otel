package webserver

import (
	"encoding/json"
	"net/http"

	"github.com.br/sk8sta13/temperatures/internal/dto"
	"github.com.br/sk8sta13/temperatures/internal/entity"
	usecase "github.com.br/sk8sta13/temperatures/use_case"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (s *WebServer) ZipCode(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := s.TemplateData.OTELTracer.Start(ctx, s.TemplateData.RequestNameOTEL+" buscando na porta 8080")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, entity.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	var requestData dto.ZipCode
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, entity.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	if requestData.ZipCode == "" {
		http.Error(w, entity.ErrZipCodeRequired.Error(), 422)
		return
	}

	if len(requestData.ZipCode) != 8 {
		http.Error(w, entity.ErrInvalidZipCode.Error(), 422)
		return
	}

	temperature, err := usecase.GetA(ctx, &requestData)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == entity.ErrCanNotFindZipcode {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	s.TemplateData.Content = temperature.City

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(temperature); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
