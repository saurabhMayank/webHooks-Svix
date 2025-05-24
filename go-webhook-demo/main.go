package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Proper initialization for Go 1.20+

	// Webhook endpoint with random failures
	e.POST("/webhook", func(c echo.Context) error {
		
		// Simulate 50% failure rate
		if r.Intn(2) == 0 { // Use the local `r` instead of global `rand`

			c.Logger().Error("Simulating a server error!")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "error: webhook failed randomly",
			})
		}

		// Process successful webhook
		var payload map[string]interface{}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "error: invalid payload",
			})
		}

		c.Logger().Info("Webhook processed successfully:", payload)
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
