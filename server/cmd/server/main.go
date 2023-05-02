package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/pkg/database"
	"server/pkg/handler"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := setupLogging(); err != nil {
		log.Fatalf("failed to setup logging: %v\n", err)
	}

	log.Println("starting server...")

	if err := database.Init(); err != nil {
		log.Fatalf("failed to initialize database: %v\n", err)
	}

	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: handler.New(),
	}

	go func() {
		err := server.ListenAndServeTLS(os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY"))
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v\n", err)
		}
	}()

	log.Printf("listening on port %v\n", server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shuting down server")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v\n", err)
	}
}

func setupLogging() error {
	logFilePath := os.Getenv("LOG_FILE")
	if logFilePath != "" {
		logFile, err := os.Create(logFilePath)
		if err != nil {
			return err
		}
		log.SetOutput(logFile)
	}
	return nil
}
