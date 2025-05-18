package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com.br/sk8sta13/temperatures/internal/dto"
	"github.com.br/sk8sta13/temperatures/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Address struct {
	Code         string `json:"cep"`
	State        string `json:"estado"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
	Region       string `json:"regiao"`
}

func Get(ctx context.Context, data *dto.ZipCode) (*dto.Temperature, error) {
	address, err := getLocal(ctx, data.ZipCode)
	if err != nil {
		return nil, err
	}

	temperature, err := getTemperature(ctx, &address.City)
	if err != nil {
		return nil, err
	}

	return temperature, nil
}

/*
	func getLocal(ctx context.Context, zipcode string) (*Address, error) {
		resp, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode))
		if err != nil {
			log.Println(err.Error())
			return nil, entity.ErrInternalServer
		}
		defer resp.Body.Close()

		var a Address
		err = json.NewDecoder(resp.Body).Decode(&a)
		if err != nil {
			log.Println(err.Error())
			return nil, entity.ErrInternalServer
		}

		if a.City == "" {
			return nil, entity.ErrCanNotFindZipcode
		}

		return &a, nil
	}
*/
func getLocal(ctx context.Context, zipcode string) (*Address, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}
	defer resp.Body.Close()

	var a Address
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	if a.City == "" {
		return nil, entity.ErrCanNotFindZipcode
	}

	return &a, nil
}

/*
	func getTemperature(ctx context.Context, region *string) (*dto.Temperature, error) {
		regionScape := url.QueryEscape(*region)
		resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?q=%s&key=%s", regionScape, os.Getenv("WEATHER_API_KEY")))
		if err != nil {
			log.Println(err.Error())
			return nil, entity.ErrInternalServer
		}
		defer resp.Body.Close()

		var d struct {
			Current dto.Temperature `json:"current"`
		}

		err = json.NewDecoder(resp.Body).Decode(&d)
		if err != nil {
			log.Println(err.Error())
			return nil, entity.ErrInternalServer
		}

		d.Current.Temp_K = d.Current.Temp_C + 273.15
		d.Current.City = *region

		return &d.Current, nil
	}
*/
func getTemperature(ctx context.Context, region *string) (*dto.Temperature, error) {
	regionScape := url.QueryEscape(*region)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?q=%s&key=%s", regionScape, os.Getenv("WEATHER_API_KEY"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}
	defer resp.Body.Close()

	var d struct {
		Current dto.Temperature `json:"current"`
	}

	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	d.Current.Temp_K = d.Current.Temp_C + 273.15
	d.Current.City = *region

	return &d.Current, nil
}
