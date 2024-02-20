package service

import (
	"context"

	"github.com/amiosamu/weather-api/internal/entity"
	"github.com/amiosamu/weather-api/internal/repo"
)

type WeatherService struct {
	weatherRepo repo.Weather
}

func (w *WeatherService) GetWeatherByCity(ctx context.Context, city string) (entity.Weather, error) {
	weather, err := w.weatherRepo.GetWeatherByCity(ctx, city)
	if err != nil {
		return entity.Weather{}, err
	}
	return weather, nil
}

func (w *WeatherService) UpdateWeather(ctx context.Context, city string) (entity.Weather, error) {
	weatherData, err := w.weatherRepo.UpdateWeather(ctx, city)
	if err != nil {
		return entity.Weather{}, err
	}
	return weatherData, nil
}

func NewWeatherService(weatherRepo repo.Weather) *WeatherService {
	return &WeatherService{weatherRepo: weatherRepo}
}
