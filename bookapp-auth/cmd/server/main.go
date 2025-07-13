package main

import (
	"log"
	"net/http"
	"time"

	"bookapp/internal/auth"
	"bookapp/internal/handlers"
	"bookapp/internal/store"
)

func main() {
	// 1) initialize in‑memory store & auth service
	userStore := store.NewInMemoryUserStore()
	authSvc := auth.NewService(userStore)

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
