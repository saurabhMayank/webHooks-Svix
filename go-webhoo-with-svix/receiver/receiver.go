package receiver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
)

// sender and receiver share the same secret key
const secretKeyReciever = "mysecretkey"

// ProcessWebhook processes the incoming webhook payload
func ProcessWebhook(c echo.Context) error {
	log.Println("Received webhook request inside the receiver")

	// Simulate random network failures in 50% of calls
	// it could be network error, it could some internal server error, it could resource is not ready
	if rand.Intn(2) == 0 { // Randomly return an error for 50% of calls
		log.Println("Simulating random failure 50% of the time")
		return echo.NewHTTPError(http.StatusInternalServerError, "Simulated network error")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	signature := c.Request().Header.Get("X-HMAC-Signature")
	if !ValidateHMAC(body, signature) {
		log.Println("Invalid HMAC signature")
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

// ValidateHMAC validates the HMAC signature
func ValidateHMAC(data []byte, signature string) bool {
	h := hmac.New(sha256.New, []byte(secretKeyReciever))
	h.Write(data)
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
