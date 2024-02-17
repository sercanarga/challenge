package main

import (
	"challenge/docs"
	"challenge/internal/durable"
	"challenge/internal/routes"
	"flag"
	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

var (
	dsn         string
	kafkaBroker string
	debug       = flag.Bool("debug", false, "debug mode")
)

func init() {
	flag.Parse()

	// setup logger
	durable.SetupLogger()

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to database
	if *debug {
		dsn = os.Getenv("DB_DSN_TEST")
	} else {
		dsn = os.Getenv("DB_DSN")
	}
	if err := durable.ConnectDB(dsn); err != nil {
		log.Fatal("Error connecting to database")
	}

	// connect to kafka
	if *debug {
		kafkaBroker = os.Getenv("KAFKA_BROKER_TEST")
	} else {
		kafkaBroker = os.Getenv("KAFKA_BROKER")
	}
	err := durable.SetupKafkaProducer(kafkaBroker)
	if err != nil {
		log.Fatal("Error connecting to kafka")
	}
}

// @version 1.0
// @description frontend service
// @BasePath /
func main() {
	if err := durable.Connection().AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
		MaxAge:        12 * time.Hour,
	}))

	// prometheus metrics
	p := ginprom.New(
		ginprom.Engine(app),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	app.Use(p.Instrument())

	// swagger routes
	docs.SwaggerInfo.Title = os.Getenv("APP_NAME")
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// routes
	app.POST("/", routes.BalanceUpdate)

	if err := app.Run(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err)
	}
}
