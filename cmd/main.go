package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mogw/micro-company/internal/auth"
	"github.com/mogw/micro-company/internal/company"
	"github.com/mogw/micro-company/internal/config"
	"github.com/mogw/micro-company/internal/db"
	"github.com/mogw/micro-company/internal/kafka"
)

func main() {
	cfg := config.LoadConfig()

	mongoClient, err := db.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	kafkaProducer := kafka.NewProducer(cfg.KafkaBroker)
	defer kafkaProducer.Close()

	companyRepo := company.NewRepository(mongoClient, "companydb", "companies")
	companyService := company.NewService(companyRepo, kafkaProducer)
	companyHandler := company.NewHandler(companyService)

	router := gin.Default()
	router.Use(auth.JWTMiddleware(cfg.JWTSecret))

	companyHandler.RegisterRoutes(router)

	log.Println("Server is running on port 8080")
	log.Fatal(router.Run(":8080"))
}
