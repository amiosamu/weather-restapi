package mongodb

import (
	"context"
	"fmt"
	"github.com/amiosamu/weather-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s", cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Addr)
	DB, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}
	return DB, err

}

func ShutdownDB(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}