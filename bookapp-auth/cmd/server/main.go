package main

import (
	"bookapp/internal/auth"
	"bookapp/internal/handlers"
	"bookapp/internal/store"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// load DSN from ENV
	dsn := os.Getenv("DATABASE_URL")
	// log.Printf("MINIO_ENDPOINT=%q", dsn)
	pgStore, err := store.NewPostgresStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	authSvc := auth.NewService(pgStore)

	// 2) wire up handlers
	router := handlers.NewRouter(authSvc)

	// 3) start server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Listening on :8080…")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
