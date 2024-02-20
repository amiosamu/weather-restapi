package service

import (
	"context"

	"github.com/amiosamu/weather-api/internal/entity"
	"github.com/amiosamu/weather-api/internal/repo"
)

type Weather interface {
	GetWeatherByCity(ctx context.Context, city string) (entity.Weather, error)
	UpdateWeather(ctx context.Context, city string)(entity.Weather, error)
}

type Services struct {
	Weather Weather
}
type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{Weather: NewWeatherService(deps.Repos.Weather)}
}


