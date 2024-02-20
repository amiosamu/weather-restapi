package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/amiosamu/weather-api/internal/entity"
	"github.com/amiosamu/weather-api/internal/repo/repoerrors"
	"github.com/amiosamu/weather-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WeatherRepo struct {
	mg *mongo.Collection
}

func (w *WeatherRepo) GetWeatherByCity(ctx context.Context, city string) (entity.Weather, error) {
	var weather entity.Weather
	filter := bson.M{"city": city}

	err := w.mg.FindOne(ctx, filter).Decode(&weather)
	if err != nil {
		fmt.Printf("Error getting weather for city %s: %s\n", city, err.Error())

		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Weather{}, repoerrors.ErrWeatherNotFound
		}

		fmt.Printf("Error getting weather for city %s: %s\n", city, err.Error())
		return entity.Weather{}, err
	}

	return weather, nil
}

func (w *WeatherRepo) UpdateWeather(ctx context.Context, city string) (entity.Weather, error) {
	apiKey := os.Getenv("WEATHER_API")
	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return entity.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entity.Weather{}, fmt.Errorf("could not request OpenWeather API: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entity.Weather{}, err
	}

	var weatherData map[string]interface{}
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return entity.Weather{}, err
	}

	temperatureKelvin := weatherData["main"].(map[string]interface{})["temp"].(float64)
	temperatureCelsius := utils.BeautifyCelsius(temperatureKelvin - 273.15)

	existingWeather, err := w.GetWeatherByCity(ctx, city)
	if err != nil && !errors.Is(err, repoerrors.ErrWeatherNotFound) {
		return entity.Weather{}, err
	}

	if existingWeather.Temperature != temperatureCelsius {
		filter := bson.M{"city": city}
		update := bson.M{"$set": bson.M{"temperature": temperatureCelsius, "last_updated": time.Now()}}

		opts := options.Update().SetUpsert(true)
		result, err := w.mg.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return entity.Weather{}, err
		}

		if result.UpsertedID != nil {
			newWeather := entity.NewWeather(city, temperatureCelsius)
			newWeather.ID = result.UpsertedID.(primitive.ObjectID)
			return newWeather, nil
		}

		return entity.Weather{City: city, Temperature: temperatureCelsius}, nil
	}

	return existingWeather, nil
}

func NewWeatherRepo(mg *mongo.Client) *WeatherRepo {
	return &WeatherRepo{mg.Database("weather-api").Collection("weather")}
}
