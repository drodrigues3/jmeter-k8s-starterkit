package main

import (
	"github.com/drodrigues3/jmeter-k8s-starterkit/config"
	"github.com/drodrigues3/jmeter-k8s-starterkit/database"
	"github.com/drodrigues3/jmeter-k8s-starterkit/handlers"
	"github.com/drodrigues3/jmeter-k8s-starterkit/log"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO

	"gorm.io/gorm"
)

// Global variable
var (
	cfg    *config.Config
	db     *gorm.DB
	err    error
	router *gin.Engine
)

func init() {

	// Load configurations
	cfg, err = config.LoadConfiguration()
	if err != nil {
		log.Panic().Msg("Error when try load configuration")
	}

	// Initialize the DB
	db, err = gorm.Open(sqlite.Open("jmeter.db"), &gorm.Config{})

	if err != nil {
		log.Panic().Err(err).Msg("Erro to connect in the DB")
	}

	// Create the scenario directory
	handlers.CreateScenarioDirectory(cfg)

	// Migrate the schema
	database.AutoMigrate(db)

	// Initialize the Gin router
	router = gin.Default()
}

func main() {

	router.Use(gin.Recovery())

	// Define the routes
	handlers.Routes(router, db, cfg)

	// Start the server
	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}
