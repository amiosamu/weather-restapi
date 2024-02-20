package weatherapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amiosamu/weather-api/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockWeatherService struct {
	GetWeatherByCityCalled bool
	UpdateWeatherCalled    bool
	Weather                entity.Weather
	Err                    error
}

func (m *MockWeatherService) GetWeatherByCity(ctx context.Context, city string) (entity.Weather, error) {
	m.GetWeatherByCityCalled = true
	return m.Weather, m.Err
}

func (m *MockWeatherService) UpdateWeather(ctx context.Context, city string) (entity.Weather, error) {
	m.UpdateWeatherCalled = true
	return m.Weather, m.Err
}

func TestGetWeather(t *testing.T) {
	router := gin.New()
	objectID, err := primitive.ObjectIDFromHex("65d49ef148642975601127a0")
	if err != nil {
		t.Fatal(err)
	}

	mockWeatherService := &MockWeatherService{
		Weather: entity.Weather{
			ID:          objectID,
			City:        "Vienna",
			Temperature: 20.0,
		},
	}

	NewWeatherRoutes(router.Group("/api/weather"), mockWeatherService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/weather/TestCity", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response getWeatherResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "65d49ef148642975601127a0", response.ID)
	assert.Equal(t, "Vienna", response.City)
	assert.Equal(t, 20.0, response.Temperature)
}
