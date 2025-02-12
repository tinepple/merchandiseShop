package main

import (
	"MerchandiseShop/internal/handler"
	"MerchandiseShop/internal/storage"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	storageRepo, err := storage.New(db)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	apiHandler := handler.New(storageRepo)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APISERVER_PORT")), apiHandler); err != nil {
		log.Fatal(err)
	}
}
