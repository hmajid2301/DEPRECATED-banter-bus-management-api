package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"banter-bus-server/src/core/database"
	"banter-bus-server/src/server"
	"banter-bus-server/src/utils/config"
)

func main() {
	appConfig := config.GetConfig()
	dbConfig := database.Config{
		Username:     appConfig.Database.Username,
		Password:     appConfig.Database.Password,
		DatabaseName: appConfig.Database.DatabaseName,
		Host:         appConfig.Database.Host,
		Port:         appConfig.Database.Port,
	}
	database.InitialiseDatabase(dbConfig)

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
