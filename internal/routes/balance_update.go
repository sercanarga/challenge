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
// @Success		200		{object}	models.APIReturn
// @Failure		400		{object}	models.APIReturn
// @Failure		500		{object}	models.APIReturn
// @Router		/		[post]
func BalanceUpdate(ctx *gin.Context) {
	var req models.EventList
	if ctx.BindJSON(&req) != nil {
		ctx.JSON(400, models.APIReturn{
			StatusCode:   400,
			Response:     "Invalid JSON format",
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	// Validate the request
	for _, event := range req.Events {
		if event.App == "" || event.Type == "" || event.Time == "" || event.Meta.User == "" || event.Wallet == "" || event.Attributes.Amount == "" || event.Attributes.Currency == "" {
			ctx.JSON(400, models.APIReturn{
				StatusCode:   400,
				Response:     "Missing required fields",
				ResponseTime: time.Now().Unix(),
			})
			return
		}
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(400, models.APIReturn{
			StatusCode:   400,
			Response:     "Failed to marshal req:" + err.Error(),
			ResponseTime: time.Now().Unix(),
		})
		return
	}
	reqString := string(reqBytes)

	message := &sarama.ProducerMessage{
		Topic: os.Getenv("KAFKA_TOPIC"),
		Value: sarama.StringEncoder(reqString),
	}

	_, _, err = durable.KafkaConnection().SendMessage(message)
	if err != nil {
		ctx.JSON(500, models.APIReturn{
			StatusCode:   500,
			Response:     "Failed to send message to Kafka:" + err.Error(),
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	ctx.Status(200)
}
