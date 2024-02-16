package main

import (
	"challenge/internal/durable"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	debug = flag.Bool("debug", false, "debug mode")
)

func init() {
	// setup logger
	durable.SetupLogger()

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to database
	if err := durable.ConnectDB(&durable.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}); err != nil {
		log.Fatal("Error connecting to database")
	}
	// durable.Connection()

}
