package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fleet-backend/config"
	"fleet-backend/internal/handler"
	"fleet-backend/internal/repository"
	"fleet-backend/internal/route"
	"fleet-backend/internal/usecase"
	"fleet-backend/internal/utils"
	"fleet-backend/pkg/db"
	"fleet-backend/pkg/mqtt"
	"fleet-backend/pkg/rabbit"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.Load()

	// DB
	dsn := utils.DSN(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSL)
	pg := utils.Must(db.Open(dsn))

	// Repo & Rabbit
	repo := repository.NewLocationPg(pg)
	rconn := utils.Must(amqp091.Dial(cfg.RabbitURL))
	pub := utils.Must(rabbit.NewPublisher(rconn, cfg.RabbitEx, cfg.RabbitExType, cfg.RabbitQueue, cfg.RabbitKey))

	// Usecase
	uc := usecase.NewFleetUsecase(repo, pub, cfg.GeofenceLat, cfg.GeofenceLon, cfg.GeofenceRad, cfg.RabbitEx, cfg.RabbitKey)

	// MQTT
	subs := mqtt.NewSubscriber(cfg.MqttURL, cfg.MqttClient, cfg.MqttTopic, uc)
	ctx, cancel := context.WithCancel(context.Background())
	if err := subs.Start(ctx); err != nil { log.Fatal(err) }

	// HTTP
	r := gin.Default()
	route.Setup(r, handler.NewVehicleHandler(uc))
	go func() { _ = r.Run(":" + cfg.AppPort) }()

	// Wait for exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	log.Println("shutting down")
}
