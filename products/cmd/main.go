package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yourname/products/internal/broker"
	"github.com/yourname/products/internal/config"
	"github.com/yourname/products/internal/db"
	"github.com/yourname/products/internal/handler"
	"github.com/yourname/products/internal/repository"
	"github.com/yourname/products/internal/service"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	time.Sleep(5 * time.Second)
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	time.Sleep(5 * time.Second)
	msgBroker, err := broker.NewRabbitMQBroker(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer msgBroker.Close()

	repo := repository.NewProductRepository(database)
	svc := service.NewProductService(repo, msgBroker)
	h := handler.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/products", func(r chi.Router) {
		r.Post("/", h.CreateProduct)
		r.Get("/", h.ListProducts)
		r.Delete("/{id}", h.DeleteProduct)
	})

	r.Handle("/metrics", promhttp.Handler())

	log.Printf("Products service starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
