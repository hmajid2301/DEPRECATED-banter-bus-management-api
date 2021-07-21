package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/api"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core/database"
)

func main() {
	retCode := mainLogic()
	os.Exit(retCode)
}

func mainLogic() int {
	logger := core.SetupLogger(os.Stdout)
	config, err := core.NewConfig()
	if err != nil {
		logger.Error(err.Error())
		return 1
	}
	core.UpdateFormatter(logger, config.App.Env)
	core.UpdateLogLevel(logger, config.App.LogLevel)

	db, err := database.NewMongoDB(logger,
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.MaxConns,
		config.DB.Timeout)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	defer db.CloseDB()

	env := &api.Env{Logger: logger, Conf: config, DB: db}
	router, err := api.Setup(env)
	if err != nil {
		logger.Errorf("Failed to load router %v.", err)
		return 1
	}

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Srv.Host, config.Srv.Port),
		Handler: router,
	}

	go terminateHandler(logger, &srv, config.DB.Timeout)

	logger.Info("The Banter Bus Management API is ready.")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("An Unexpected error while serving HTTP %v.", err)
		return 1
	}

	logger.Info("The Banter Bus Management API has been gracefully terminated.")
	return 0
}

// terminateHandler waits for SIGINT or SIGTERM signals and does a graceful shutdown of the HTTP server
// Wait for interrupt signal to gracefully shutdown the server with
// a timeout of 5 seconds.
// kill (no param) default send syscall.SIGTERM
// kill -2 is syscall.SIGINT
// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
func terminateHandler(logger *log.Logger, srv *http.Server, timeout int) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down HTTP server.")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Unexpected error while shutting down server %s", err)
	}
}
