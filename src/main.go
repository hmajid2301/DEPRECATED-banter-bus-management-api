package main

import (
	// Log items to the terminal
	"banter-bus-server/src/controllers"
	"banter-bus-server/src/database"
	"log"

	// Import gin for route definition
	"github.com/gin-gonic/gin"
	// Import godotenv for .env variables
	"github.com/joho/godotenv"
	// Import our app controllers
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// init gets called before the main function
func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

func main() {
	config := database.DatabaseConfig{Username: "banterbus", Password: "banterbus", DatabaseName: "banterbus", Host: "banter-bus-database", Port: "27017"}
	// Init gin router
	router := gin.Default()
	database.InitialiseDatabase(config)

	// Its great to version your API's
	v1 := router.Group("/api/v1")
	{
		// Define the hello controller
		hello := new(controllers.HelloWorldController)
		// Define a GET request to call the Default
		// method in controllers/hello.go
		v1.GET("/hello", hello.Default)
	}

	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		// In gin this is how you return a JSON response
		c.JSON(404, gin.H{"message": "Not found"})
	})

	// Init our server
	router.Run(":8080")
}

// Config is the data type for the expected config file.
type Config struct {
	Database struct {
		Host     string `yaml:"host" env:"BANTER_BUS_DB_HOST" env-default:"database"`
		Port     string `yaml:"port" env:"BANTER_BUS_DB_PORT" env-default:"27017"`
		Name     string `yaml:"name" env:"BANTER_BUS_DB_NAME" env-default:"banterbus"`
		User     string `yaml:"user" env:"BANTER_BUS_DB_USER"`
		Password string `yaml:"password" env:"BANTER_BUS_DB_PASSWORD"`
	} `yaml:"database"`
}

func config() *Config {
	var cfg Config

	path, exists := os.LookupEnv("BANTER_BUS_CONFIG_PATH")
	var configPath = "config.yaml"
	if exists {
		configPath = path
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return &cfg
}
