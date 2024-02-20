package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amiosamu/weather-api/config"
	weatherapi "github.com/amiosamu/weather-api/internal/controller/http/weather-api"
	"github.com/amiosamu/weather-api/internal/repo"
	"github.com/amiosamu/weather-api/internal/service"
	"github.com/amiosamu/weather-api/pkg/httpserver"
	"github.com/amiosamu/weather-api/pkg/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(path string) {
	cfg, err := config.NewConfig(path)
	if err != nil {
		slog.Error("error reading config: %w", err)
	}
	slog.Info("setting up MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	mg, err := mongodb.InitMongo(ctx, cfg)
	if err != nil {
		log.Fatalf("failed connecting to Mongo: %v", err)
	}
	ctxDB, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := mg.Ping(ctxDB, nil); err != nil {
		log.Fatalf("error pinging Mongo: %v", err)
	}
	defer func(ctx context.Context, client *mongo.Client) {
		err := mongodb.ShutdownDB(ctx, client)
		if err != nil {
			log.Fatalf("error shuttding down Mongo: %v", err)
		}
	}(ctxDB, mg)

	slog.Info("initializing repositories...")

	repository := repo.NewRepositories(mg)

	slog.Info("initializing service dependencies...")

	dependencies := service.ServicesDependencies{

		Repos: repository,
	}

	services := service.NewServices(dependencies)

	slog.Info("initializing handlers and routes...")

	handler := gin.New()
	weatherapi.NewRouter(handler, services)

	slog.Info("starting http server...")
	slog.Debug("server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	slog.Info("configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		slog.Error("app - Run - httpServer.Notify: " + err.Error())
	}

	slog.Info("shutting down...")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("app - Run - httpServer.Notify: " + err.Error())
	}
}
