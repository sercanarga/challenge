package main

import (
	"challenge/docs"
	"challenge/internal/durable"
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
}

// @version 1.0
// @description frontend service
// @BasePath /
func main() {
	flag.Parse()

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
		AllowMethods:  []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
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

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
