package receiver

import (
	"encoding/json"
	"go-webhoo-with-svix/configs"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	svix "github.com/svix/svix-webhooks/go"
)

// sender and receiver share the same secret key
const secretKeyReciever = "mysecretkey"

// ProcessWebhook processes the incoming webhook payload
func ProcessWebhook(c echo.Context) error {
	log.Println("Received webhook request inside the receiver")

	// Simulate random network failures in 50% of calls
	// it could be network error, it could some internal server error, it could resource is not ready
	// Create a new random number generator with a custom source
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	if rng.Intn(2) == 0 { // Randomly return an error for 50% of calls
		log.Println("Simulating random failure 50% of the time")
		return echo.NewHTTPError(http.StatusInternalServerError, "Simulated network error")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if !ValidateSvixSignature(body, c.Request().Header) {
		log.Println("Invalid svix signature")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("Failed to parse payload:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	log.Println("Webhook processed successfully:", payload)
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// ValidateSvixSignature validates the Svix signature
func ValidateSvixSignature(payload []byte, headers http.Header) bool {
	svixSigningKey := configs.GetSvixSigningKey()
	wh, err := svix.NewWebhook(svixSigningKey)
	if err != nil {
		log.Println("Failed to initialize Svix webhook verifier:", err)
		return false
	}

	signature := headers.Get("Svix-Signature")
	if signature == "" {
		log.Println("Missing Svix-Signature header")
		return false
	}

	err = wh.Verify(payload, headers)
	if err != nil {
		log.Println("Invalid Svix signature:", err)
		return false
	}

	return true
}
