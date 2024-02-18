package routes

import (
	"bytes"
	"challenge/internal/models"
	"challenge/internal/routes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBalanceUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid request", func(t *testing.T) {
		router := gin.Default()
		router.POST("/")

		eventList := models.EventList{
			Events: []models.Event{
				{
					App:  "01HPMTX8916FF4ABFBDESX1AGH",
					Type: "BALANCE_INCREASE",
					Time: "2024-02-12T11:50:40.280Z",
					Meta: models.Meta{
						User: "01HPMV114ZE7Z54M6XV8H4EEMB",
					},
					Wallet: "01HPMV01XPAXCG242W7SZWD0S5",
					Attributes: models.Attributes{
						Amount:   "33.20",
						Currency: "TRY",
					},
				},
			},
		}

		jsonData, err := json.Marshal(eventList)
		if err != nil {
			log.Fatalf("JSON Marshaling failed: %s", err)
		}

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Request creation failed: %s", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		router := gin.Default()
		router.POST("/", routes.BalanceUpdate)

		eventList := models.EventList{
			Events: []models.Event{}, // Events alanı boş
		}

		jsonData, err := json.Marshal(eventList)
		if err != nil {
			log.Fatalf("JSON Marshaling failed: %s", err)
		}

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Request creation failed: %s", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.NotEqual(t, http.StatusOK, resp.Code)
	})
}
