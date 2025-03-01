package main

import (
	"context"
	"log"
	"mucb_be/internal/app"
	"mucb_be/internal/delivery/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting application...")

	deps, cfg := app.Start()

	router := gin.Default()
	http.SetupRouter(router, cfg, deps)

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Server running on port 8080")
		serverErrors <- router.Run(":8080")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-sigChan:
		log.Printf("Shutting down gracefully... Reason: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		app.Stop(ctx, deps)

		log.Println("Server shutdown complete")
	}
}
