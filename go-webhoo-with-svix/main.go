package main

import (
	"log"
	"net/http"

	"go-webhoo-with-svix/receiver" // Import the receiver package

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Define routes
	e.POST("/webhook", func(c echo.Context) error {
		log.Println("Received webhook request") // Debug log
		return receiver.ProcessWebhook(c)
	})

	// Start server
	log.Println("Starting receiver server on :8080")
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
