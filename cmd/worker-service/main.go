package main

import (
	"challenge/internal/durable"
	"challenge/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	Broker      string
	Topic       string
	WorkerCount int
	BatchSize   int
	BatchTime   time.Duration
}

var config Config
var consumer sarama.Consumer
var err error

func init() {
	// setup logger
	durable.SetupLogger()

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to database
	if err := durable.ConnectDB(os.Getenv("DB_DSN")); err != nil {
		log.Fatal("Error connecting to database")
	}

	// create a new Kafka consumer
	config = Config{
		Broker:      os.Getenv("KAFKA_BROKER"),
		Topic:       os.Getenv("KAFKA_TOPIC"),
		WorkerCount: 5,
		BatchSize:   100,
		BatchTime:   2 * time.Second,
	}
	consumer, err = durable.SetupKafkaConsumer(config.Broker)
	if err != nil {
		log.Fatalf("Error creating consumer: %v\n", err)
	}
}

func main() {
	defer func(consumer sarama.Consumer) {
		err := durable.CloseKafkaConsumer(consumer)
		if err != nil {
			log.Fatalf("Error closing consumer: %v\n", err)
		}
	}(consumer)

	// get the list of partitions for topic
	partitions, err := consumer.Partitions(config.Topic)
	if err != nil {
		log.Fatalf("Error retrieving partitions: %v\n", err)
	}

	// create wait group for workers
	var wg sync.WaitGroup
	wg.Add(config.WorkerCount)

	// create buffered channel for messages
	msgs := make(chan *sarama.ConsumerMessage, config.BatchSize*config.WorkerCount)

	// start workers
	for i := 0; i < config.WorkerCount; i++ {
		go startWorker(&wg, msgs)
	}

	// listen for interrupt signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// start partition consumers
	for _, partition := range partitions {
		startPartitionConsumer(consumer, partition, msgs, signals)
	}

	// wait all workers to finish
	wg.Wait()
}

// starts a new worker
func startWorker(wg *sync.WaitGroup, msgs chan *sarama.ConsumerMessage) {
	defer wg.Done()
	batch := make([]*sarama.ConsumerMessage, 0, config.BatchSize)
	timer := time.NewTimer(config.BatchTime)
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				processBatch(batch)
				return
			}
			batch = append(batch, msg)
			if len(batch) >= config.BatchSize {
				processBatch(batch)
				batch = batch[:0]
				timer.Reset(config.BatchTime)
			}
		case <-timer.C:
			processBatch(batch)
			batch = batch[:0]
			timer.Reset(config.BatchTime)
		}
	}
}

// Starts a new partition consumer
func startPartitionConsumer(consumer sarama.Consumer, partition int32, msgs chan *sarama.ConsumerMessage, signals chan os.Signal) {
	partitionConsumer, err := consumer.ConsumePartition(config.Topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v\n", err)
	}

	go func(pc sarama.PartitionConsumer) {
		for {
			select {
			case msg := <-pc.Messages():
				msgs <- msg
			case err := <-pc.Errors():
				log.Printf("Error: %v\n", err) // log error
			case <-signals:
				fmt.Println("Interrupt signal received. Shutting down...")
				close(msgs)
				err := pc.Close()
				if err != nil {
					return
				}
				return
			}
		}
	}(partitionConsumer)
}

func processBatch(batch []*sarama.ConsumerMessage) {
	for _, msg := range batch {
		var event models.Event
		var user models.Users
		var wallet models.Wallet
		var balance models.Balance

		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		// check if user exists
		userResult := durable.Connection().First(&user, "id = ? AND deleted_at IS NULL", event.Meta.User)
		if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
			log.Println("User not found or deleted")
			continue
		}

		// check if wallet exists
		walletResult := durable.Connection().First(&wallet, "id = ? AND user_id = ?", event.Wallet, event.Meta.User)
		if errors.Is(walletResult.Error, gorm.ErrRecordNotFound) {
			log.Println("Wallet not found or wallet not belongs to user")
			continue
		}

		eventAmount, err := strconv.ParseFloat(event.Attributes.Amount, 64)
		if err != nil {
			log.Printf("Error parsing amount: %v\n", err)
			continue
		}

		// check if balance exists and update or create
		BalanceResult := durable.Connection().First(&balance, "wallet_id = ? AND currency = ?", event.Wallet, event.Attributes.Currency)
		if errors.Is(BalanceResult.Error, gorm.ErrRecordNotFound) {
			// wallet does not exist, create a new one
			balance = models.Balance{
				WalletId: event.Wallet,
				Currency: event.Attributes.Currency,
				Amount:   eventAmount,
			}
		} else {
			switch event.Type {
			case "BALANCE_INCREASE":
				balance.Amount += eventAmount
			case "BALANCE_DECREASE":
				balance.Amount -= eventAmount
			}
		}
		balance.LastUpdate = time.Now()

		result := durable.Connection().Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "currency"}, {Name: "wallet_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"amount", "last_update"}),
		}).Create(&balance)

		if result.Error != nil {
			log.Printf("Error updating or creating wallet: %v\n", result.Error)
			continue
		}
	}
}
