package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server"
	"banter-bus-server/src/utils/config"
	logger "banter-bus-server/src/utils/log"
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
	logger.FormatLogger()

	router, err := server.NewRouter()
	if err != nil {
		log.Fatal("Failed to load router", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}
