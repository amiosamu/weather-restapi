package utils

import "math"

func BeautifyCelsius(temp float64) float64 {
	return math.Round(temp*100) / 100
}
