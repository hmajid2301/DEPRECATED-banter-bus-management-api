package repository

import (
	log "github.com/sirupsen/logrus"
)

// NewRepository creates a new datastore object, used to store data (persistently).
func NewRepository(
	logger *log.Logger,
	host string,
	port int,
	username string,
	password string,
	databaseName string,
	maxConns int,
	timeout int,
) (*MongoDB, error) {
	db, err := NewMongoDB(
		logger,
		host,
		port,
		username,
		password,
		databaseName,
		maxConns,
		timeout)
	return db, err
}
