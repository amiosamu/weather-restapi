package weatherapi

import (
	"net/http"

	"github.com/amiosamu/weather-api/internal/service"
	"github.com/gin-gonic/gin"
)

type weatherRoutes struct {
	weatherService service.Weather
}

func NewWeatherRoutes(c *gin.RouterGroup, weatherService service.Weather) {
	r := &weatherRoutes{
		weatherService: weatherService,
	}

	c.PUT("/", r.update)
	c.GET("/:city", r.get)

}

type GetWeatherRequest struct {
	City string `json:"city"`
}

type getWeatherResponse struct {
	ID          string  `json:"id"`
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Code        int     `json:"code"`
}

type UpdateWeatherRequest struct {
	City string `json:"city"`
}

type updateStatusResponse struct {
	ID          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Code        int     `json:"code"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func (r *weatherRoutes) get(ctx *gin.Context) {

	city := ctx.Param("city")
	weather, err := r.weatherService.GetWeatherByCity(ctx, city)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, statusResponse{"Could not get weather."})
		return

	}

	resp := getWeatherResponse{
		ID:          weather.ID.Hex(),
		City:        weather.City,
		Temperature: weather.Temperature,
		Code:        http.StatusOK,
	}

	ctx.JSON(http.StatusOK, resp)

}

func (r *weatherRoutes) update(ctx *gin.Context) {
	var request UpdateWeatherRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weather, err := r.weatherService.UpdateWeather(ctx, request.City)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, statusResponse{"Could not update weather."})
		return
	}

	resp := updateStatusResponse{
		ID:          weather.ID.Hex(),
		Temperature: weather.Temperature,
		Code:        http.StatusOK,
	}

	ctx.JSON(http.StatusOK, resp)
}
