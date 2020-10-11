package main

import (
	"log"
	"net/http"

	"banter-bus-server/src/core/config"
	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server"
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
	database.InitialiseDatabase(dbConfig)

	router, err := server.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	srv.ListenAndServe()
}
