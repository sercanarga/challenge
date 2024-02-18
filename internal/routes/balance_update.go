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

	successfulEvents := make([]models.Event, 0)
	unsuccessfulEvents := make([]models.Event, 0)

	users := make(map[string]models.Users)
	wallets := make(map[string]models.Wallet)
	messages := make([]*sarama.ProducerMessage, 0)

	for _, event := range req.Events {
		var user models.Users
		var wallet models.Wallet

		// Validate the request
		if event.App == "" ||
			event.Type == "" ||
			event.Time == "" ||
			event.Meta.User == "" ||
			event.Wallet == "" ||
			event.Attributes.Amount == "" ||
			event.Attributes.Currency == "" {
			unsuccessfulEvents = append(unsuccessfulEvents, event)
			continue
		}

		// === Database checks disabled ===
		// Check if user and wallet exist
		//if _, ok := users[event.Meta.User]; !ok {
		//	userResult := durable.Connection().First(&user, "id = ?", event.Meta.User)
		//	if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		//		unsuccessfulEvents = append(unsuccessfulEvents, event)
		//		continue
		//	}
		//	users[event.Meta.User] = user
		//}
		//if _, ok := wallets[event.Wallet]; !ok {
		//	walletResult := durable.Connection().First(&wallet, "id = ? AND user_id = ?", event.Wallet, event.Meta.User)
		//	if errors.Is(walletResult.Error, gorm.ErrRecordNotFound) {
		//		unsuccessfulEvents = append(unsuccessfulEvents, event)
		//		continue
		//	}
		//	wallets[event.Wallet] = wallet
		//}

		users[event.Meta.User] = user
		wallets[event.Wallet] = wallet
		// === Database checks disabled ===

		// Marshal the event
		eventJson, err := json.Marshal(event)
		if err != nil {
			ctx.JSON(400, models.APIReturn{
				StatusCode:   400,
				Response:     "Failed to marshal req:" + err.Error(),
				ResponseTime: time.Now().Unix(),
			})
			return
		}

		// Create Kafka message
		message := &sarama.ProducerMessage{
			Topic: os.Getenv("KAFKA_TOPIC"),
			Value: sarama.StringEncoder(eventJson),
		}
		messages = append(messages, message)
		successfulEvents = append(successfulEvents, event)
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
	if len(unsuccessfulEvents) > 0 {
		ctx.JSON(207, models.APIEventReturn{
			StatusCode:   207,
			Success:      successfulEvents,
			Unsuccess:    unsuccessfulEvents,
			ResponseTime: time.Now().Unix(),
		})
		return
	}

	ctx.JSON(202, models.APIEventReturn{
		StatusCode:   202,
		Success:      successfulEvents,
		Unsuccess:    unsuccessfulEvents,
		ResponseTime: time.Now().Unix(),
	})
}
