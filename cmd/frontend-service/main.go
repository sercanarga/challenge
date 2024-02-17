package main

import (
	"challenge/internal/durable"
	"flag"
	"fmt"
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
	if err := durable.ConnectDB(os.Getenv("DB_DSN")); err != nil {
		log.Fatal("Error connecting to database")
	}
	// durable.Connection()

}

func main() {
	durable.Connection().AutoMigrate()
	fmt.Println("Connected to database")
}
