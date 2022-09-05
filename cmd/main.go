package main

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/rusystem/cache"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/repository"
	"github.com/rusystem/notes-app/internal/server"
	"github.com/rusystem/notes-app/internal/service"
	grpc_client "github.com/rusystem/notes-app/internal/transport/grpc"
	"github.com/rusystem/notes-app/internal/transport/rest"
	"github.com/rusystem/notes-app/pkg/database"
	"github.com/sirupsen/logrus"
	"log"

	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.ErrorLevel)
}

// @title Note app API
// @version 1.0
// @description API server for Note app

// @contact.name Dmitry Mikhaylov
// @contact.email ru.system.ru@gmail.com

// @host localhost:8080
// @BasePath /

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := database.NewPostgresConnection(database.PSQLConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			logrus.Fatal(err)
		}
	}(db)

	rdb := database.NewRedisClient(database.RedisConnectionInfo{
		Host:     cfg.RDB.Host,
		Port:     cfg.RDB.Port,
		Password: cfg.RDB.Password,
	})
	defer func(rdb *redis.Client) {
		if err := rdb.Close(); err != nil {
			logrus.Fatal(err)
		}
	}(rdb)

	c := cache.New()

	logsClient, err := grpc_client.NewClient(cfg.Grpc.Host, cfg.Grpc.Port)
	if err != nil {
		log.Fatal(err)
	}

	noteRepo := repository.NewRepository(db, rdb)
	noteService := service.NewService(cfg, c, noteRepo, logsClient)
	handler := rest.NewHandler(noteService)

	srv := server.New(cfg, handler.InitRoutes())
	go func() {
		if err := srv.Run(); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Note-app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Note-app stopped")

	if err := srv.Stop(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}
}
