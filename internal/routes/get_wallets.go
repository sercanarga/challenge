package routes

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

// GetWallets
// @Summary		Returns a list of wallets.
// @Produce		json
// @Success		200		{object}	models.APIReturn
// @Failure		400		{object}	models.APIReturn
// @Failure		500		{object}	models.APIReturn
// @Router		/		[get]
func GetWallets(ctx *gin.Context) {
	/*
		@todo: change to use the database to get the wallets.
			   currently returning the messages from kafka
	*/

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Fatalln("Failed to start consumer:", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalln("Failed to start partition consumer:", err)
	}

	messages := make([]string, 0)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			messages = append(messages, string(msg.Value))
			if msg.Offset == partitionConsumer.HighWaterMarkOffset()-1 {
				ctx.JSON(200, gin.H{
					"messages": messages,
				})
				return
			}
		case err := <-partitionConsumer.Errors():
			log.Fatalln("Failed to consume partition:", err)
		}
	}
}
