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
	response := make([]models.APIEventReturn, 0)

	for _, event := range req.Events {
		var r = models.APIEventReturn{
			Data: event,
		}

		// Validate the request
		if r.Data.App == "" || r.Data.Type == "" || r.Data.Time == "" || r.Data.Meta.User == "" || r.Data.Wallet == "" || r.Data.Attributes.Amount == "" || r.Data.Attributes.Currency == "" {
			r.Result.StatusCode = 400
			r.Result.ErrorDetails = "Invalid event data"
		} else {
			// Marshal the event
			eventJson, err := json.Marshal(event)
			if err != nil {
				r.Result.StatusCode = 400
				r.Result.ErrorDetails = "Failed to marshal req:" + err.Error()
			} else {
				// Create Kafka message
				message := &sarama.ProducerMessage{
					Topic: os.Getenv("KAFKA_TOPIC"),
					Value: sarama.StringEncoder(eventJson),
				}
				messages = append(messages, message)

				r.Result.StatusCode = 202
			}
		}

		response = append(response, r)
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
	ctx.JSON(202, response)
}
