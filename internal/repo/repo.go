package repo

import (
	"context"

	"github.com/amiosamu/weather-api/internal/entity"
	"github.com/amiosamu/weather-api/internal/repo/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Weather interface {
	GetWeatherByCity(ctx context.Context, city string) (entity.Weather, error)
	UpdateWeather(ctx context.Context, city string) (entity.Weather, error)
}

type Repositories struct {
	Weather
}

func NewRepositories(mg *mongo.Client) *Repositories {
	return &Repositories{
		Weather: mongodb.NewWeatherRepo(mg),
	}
}
