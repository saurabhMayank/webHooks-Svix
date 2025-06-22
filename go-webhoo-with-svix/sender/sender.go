package sender

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	webhookURL   = "http://localhost:8080/webhook"
	maxRetries   = 5
	initialDelay = time.Second
	// secretKeySender = "mysecretkeyNew" // key changed, it should break now
	secretKeySender = "mysecretkey" // key changed, it should break now
)

// SecurePayload encrypts and signs the payload using HMAC
func SecurePayload(payload map[string]interface{}) (string, string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	// Generate HMAC signature
	h := hmac.New(sha256.New, []byte(secretKeySender))
	h.Write(data)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return string(data), signature, nil
}

// PrintCurlCommand prints the equivalent curl command for the API request
func PrintCurlCommand(payload string, signature string) {
	fmt.Println("Equivalent cURL command:")
	fmt.Printf(`curl -X POST %s -H "Content-Type: application/json" -H "X-HMAC-Signature: %s" -d '%s'\n`, webhookURL, signature, payload)
}

// SendWebhook sends the webhook with retries
func SendWebhook(payload map[string]interface{}) error {
	data, signature, err := SecurePayload(payload)
	if err != nil {
		return fmt.Errorf("failed to secure payload: %v", err)
	}

	// Validate URL before any attempts
	if _, err := url.ParseRequestURI(webhookURL); err != nil {
		return fmt.Errorf("invalid URL: %w", err) // Abort immediately
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-HMAC-Signature", signature)

		// Print the cURL command
		PrintCurlCommand(data, signature)

		resp, err := http.DefaultClient.Do(req)

		if err != nil { // some network glitch
			log.Printf("Attempt %d: Failed with HTTP status: %v", attempt+1, resp.Status)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusUnauthorized {
				log.Println("Unauthorized: Invalid HMAC signature. Please check the secret key.")
				return fmt.Errorf("unauthorized: invalid HMAC signature")
			} else if resp.StatusCode == http.StatusBadRequest {
				log.Println("Bad Request: The payload might be invalid or malformed.")
				return fmt.Errorf("bad request: the payload might be invalid or malformed")
			} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				// add a new line for better readability
				log.Println()
				log.Println("Webhook sent successfully!")
				return nil
			} else {
				log.Printf("Attempt %d: Failed with HTTP status: %v", attempt+1, resp.Status)
			}
		}

		// Exponential backoff
		delay := time.Duration(math.Pow(2, float64(attempt))) * initialDelay
		log.Printf("Waiting %v before retry...", delay)
		time.Sleep(delay)
	}
	return fmt.Errorf("all %d retries failed", maxRetries)
}

func main() {
	// Check if payload is passed as an argument
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run sender.go '{\"event\":\"test_event\",\"data\":\"sample_data\"}'")
	}

	// Parse the JSON payload from the command line argument
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(os.Args[1]), &payload)
	if err != nil {
		log.Fatalf("Invalid JSON payload: %v", err)
	}

	// Send the webhook
	if err := SendWebhook(payload); err != nil {
		log.Fatalf("Failed to send webhook: %v", err)
	}
}
