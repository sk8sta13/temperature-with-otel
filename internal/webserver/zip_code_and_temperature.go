package webserver

import (
	"encoding/json"
	"net/http"

	"github.com.br/sk8sta13/temperatures/internal/dto"
	"github.com.br/sk8sta13/temperatures/internal/entity"
	"github.com.br/sk8sta13/temperatures/internal/validators"
	usecase "github.com.br/sk8sta13/temperatures/use_case"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (s *WebServer) ZipCodeAndTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := s.TemplateData.OTELTracer.Start(ctx, s.TemplateData.RequestNameOTEL+" chamada externa")
	defer span.End()

	var requestData dto.ZipCode
	requestData.ZipCode = r.URL.Query().Get("zipcode")

	if requestData.ZipCode == "" {
		http.Error(w, entity.ErrZipCodeRequired.Error(), 422)
		return
	}

	if !validators.IsValidZipCode(requestData.ZipCode) {
		http.Error(w, entity.ErrInvalidZipCode.Error(), 422)
		return
	}

	temperatures, err := usecase.Get(ctx, &requestData)
	s.TemplateData.Content = temperatures.City
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == entity.ErrCanNotFindZipcode {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(temperatures); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
