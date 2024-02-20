package weatherapi

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amiosamu/weather-api/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, services *service.Services) *gin.Engine {

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf(`{"time":"%s", "method":"%s","uri":"%s", "status":%d,"error":"%s"}`,
				param.TimeStamp.Format(time.RFC3339Nano),
				param.Method,
				param.Path,
				param.StatusCode,
				param.ErrorMessage,
			)
		},
		Output: setLogsFile(),
	}))
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	todoList := router.Group("/api")
	{
		NewWeatherRoutes(todoList.Group("/weather"), services.Weather)
	}

	return router
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
