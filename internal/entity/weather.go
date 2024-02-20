package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Weather struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	City        string             `bson:"city" json:"city"`
	Temperature float64            `bson:"temperature" json:"temperature"`
	LastUpdated primitive.DateTime `bson:"last_updated" json:"last_updated"`
}

func NewWeather(city string, temperature float64) Weather {
	return Weather{
		City:        city,
		Temperature: temperature,
		LastUpdated: primitive.NewDateTimeFromTime(time.Now()),
	}
}
