package main

import "github.com/amiosamu/weather-api/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
