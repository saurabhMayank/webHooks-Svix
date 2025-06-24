package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"go-webhoo-with-svix/configs"
	"log"
	"os"

	svix "github.com/svix/svix-webhooks/go"
)

// SendWebhook sends the webhook with retries
func SendWebhook(payload map[string]interface{}) error {
	// replacing svix with manual code of sending webhook and doing retries
	// when sending webhook request manually need to take care of retries and encrypting the body
	// which svix takes care of

	// svix implementation
	svixApiKey := configs.GetSvixKey()
	svixAppId := configs.GetSvixAppID()

	if svixApiKey == "" || svixAppId == "" {
		return fmt.Errorf("svix_key or svix_app_id is missing in the configuration")
	}

	// Initialize the Svix client
	client, err := svix.New(svixApiKey, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Svix client: %v", err)
	}

	// Define the webhook payload
	eventType := "example.event" // Replace with your event type
	// payloadBytes, err := json.Marshal(payload)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal payload: %v", err)
	// }

	// Send the webhook using Svix
	_, err = client.Message.Create(context.Background(), svixAppId, svix.MessageIn{
		EventType: eventType,
		Payload:   payload,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to send webhook via Svix: %v", err)
	}

	log.Println("Webhook sent successfully via Svix!")
	return nil
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
