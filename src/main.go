package main

import (
	"banter-bus-server/src/config"
	"banter-bus-server/src/controllers"
	"banter-bus-server/src/database"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.GetConfig()
	dbConfig := database.DatabaseConfig{
		Username:     config.Database.Username,
		Password:     config.Database.Password,
		DatabaseName: config.Database.DatabaseName,
		Host:         config.Database.Host,
		Port:         config.Database.Port,
	}
	router := gin.Default()
	database.InitialiseDatabase(dbConfig)

	v1 := router.Group("/api/v1")
	{
		hello := new(controllers.HelloWorldController)
		v1.GET("/hello", hello.Default)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	router.Run(":8080")
}
