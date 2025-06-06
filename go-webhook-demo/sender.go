package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"
)

// dummy commit
const (
	maxRetries   = 3
	initialDelay = 1 * time.Second
	webhookURL   = "http://localhost:8080/webhook" // Define URL as a constant
)

func sendWebhookWithRetries(payload []byte) error {
	// Validate URL before any attempts
	if _, err := url.ParseRequestURI(webhookURL); err != nil {
		return fmt.Errorf("invalid URL: %w", err) // Abort immediately
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Attempt %d: Connection error: %v (will retry)", attempt+1, err)

		} else {

			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body) // ← Add this for empty responses

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				log.Printf("Attempt %d: Success!", attempt+1)
				return nil
			}

			log.Printf("Attempt %d: Failed with HTTP status: %s", attempt+1, resp.Status)
		}

		delay := time.Duration(math.Pow(2, float64(attempt))) * initialDelay
		log.Printf("Waiting %v before retry...", delay)
		time.Sleep(delay)

	}
	return fmt.Errorf("all %d retries failed", maxRetries)
}

func main() {
	payload := []byte(`{"event": "user.signup", "user_id": 123}`)

	if err := sendWebhookWithRetries(payload); err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			log.Fatalf("Configuration error (fix URL!): %v", urlErr)
		} else {
			log.Fatalf("Delivery failed after retries: %v", err)
		}
	} else {
		log.Println("Webhook delivered successfully!")
	}
}
