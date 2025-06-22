package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/capernix/gohttpx/handlers"
	"github.com/capernix/gohttpx/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", handlers.CreateUser)
	mux.HandleFunc("GET /users", handlers.ListUsers)
	mux.HandleFunc("GET /users/{id}", handlers.GetUser)
	mux.HandleFunc("DELETE /users/{id}", handlers.DeleteUser)

	mux.HandleFunc("POST /notes", handlers.CreateNote)
	mux.HandleFunc("GET /notes", handlers.ListNotes)
	mux.HandleFunc("GET /notes/{id}", handlers.GetNote)
	mux.HandleFunc("DELETE /notes/{id}", handlers.DeleteNote)

	handler := middleware.Chain(
		middleware.Logger,
		middleware.Recovery,
	)(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		fmt.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
