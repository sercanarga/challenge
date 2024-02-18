package routes

import (
	"challenge/internal/durable"
	"challenge/internal/models"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

// BalanceUpdate
// @Summary		Increases or decrease the balance of the respective currency in a wallet.
// @Accept		json
// @Produce		json
// @Param		_		body		models.EventList	true "raw"
// @Success		202		{object}	models.APIEventReturn
// @Success		207		{object}	models.APIEventReturn
// @Failure		400		{object}	models.APIReturn
// @Failure		500		{object}	models.APIReturn
// @Router		/		[post]
func BalanceUpdate(ctx *gin.Context) {
	var req models.EventList

	if ctx.BindJSON(&req) != nil || req.Events == nil {
		ctx.JSON(400, models.APIReturn{
			StatusCode:   400,
			Response:     "Invalid JSON format",
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	messages := make([]*sarama.ProducerMessage, 0)

	for i, event := range req.Events {
		// Validate the request
		if event.App == "" ||
			event.Type == "" ||
			event.Time == "" ||
			event.Meta.User == "" ||
			event.Wallet == "" ||
			event.Attributes.Amount == "" ||
			event.Attributes.Currency == "" {
			req.Events[i].Response.StatusCode = 400
			req.Events[i].Response.ErrorDetails = "Invalid event data"
			continue
		}

		// Marshal the event
		eventJson, err := json.Marshal(event)
		if err != nil {
			req.Events[i].Response.StatusCode = 400
			req.Events[i].Response.ErrorDetails = "Failed to marshal req:" + err.Error()
			continue
		}

		// Create Kafka message
		message := &sarama.ProducerMessage{
			Topic: os.Getenv("KAFKA_TOPIC"),
			Value: sarama.StringEncoder(eventJson),
		}
		messages = append(messages, message)
		req.Events[i].Response.StatusCode = 202
	}

	// Send messages to Kafka
	err := durable.KafkaConnection().SendMessages(messages)
	if err != nil {
		ctx.JSON(500, models.APIReturn{
			StatusCode:   500,
			Response:     "Failed to send messages to Kafka:" + err.Error(),
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	// Return the response
	ctx.JSON(202, models.EventList{
		Events: req.Events,
	})
}
