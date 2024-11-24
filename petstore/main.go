package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/evanlindsey/go-micro/petstore/api"
	"github.com/evanlindsey/go-micro/petstore/ent"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbName, dbUser, dbPass)

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	server := api.NewServer(client)
	h := api.NewStrictHandler(server, nil)
	r := chi.NewRouter()
	api.HandlerFromMux(h, r)
	s := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	log.Println("Starting server on http://localhost:8080")
	log.Fatal(s.ListenAndServe())
}
