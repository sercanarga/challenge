package main

import (
	"challenge/docs"
	"challenge/internal/durable"
	"challenge/internal/models"
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
	debug = flag.Bool("debug", false, "debug mode")
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
	if err := durable.ConnectDB(os.Getenv("DB_DSN")); err != nil {
		log.Fatal("Error connecting to database")
	}

	// connect to kafka
	err := durable.SetupKafkaProducer(os.Getenv("KAFKA_BROKER"))
	if err != nil {
		log.Fatal("Error connecting to kafka")
	}
}

// @version 1.0
// @description frontend service
// @BasePath /
func main() {
	if err := durable.Connection().AutoMigrate(
		&models.Users{},
		&models.Wallet{},
		&models.Balance{},
	); err != nil {
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
	app.GET("/doc", func(c *gin.Context) {
		c.Redirect(301, "/docs/index.html")
	})

	// routes
	app.POST("/", routes.BalanceUpdate)
	app.GET("/", routes.GetWallets)

	if err := app.Run(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err)
	}
}
