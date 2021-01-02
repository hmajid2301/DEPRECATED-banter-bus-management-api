package core

import (
	"fmt"
	"net"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the data type for the expected config file.
type Config struct {
	App struct {
		Environment string `yaml:"environment" env:"BANTER_BUS_ENVIRONMENT" env-default:"production"`
		LogLevel    string `yaml:"logLevel" env:"BANTER_BUS_LOG_LEVEL" env-default:"debug"`
	} `yaml:"app"`
	Webserver struct {
		Host string `yaml:"host" env:"BANTER_BUS_WEBSERVER_HOST" env-default:"0.0.0.0"`
		Port int    `yaml:"port" env:"BANTER_BUS_WEBSERVER_PORT" env-default:"8080"`
	} `yaml:"webserver"`
	Database struct {
		Host         string `yaml:"host" env:"BANTER_BUS_DB_HOST" env-default:"banter-bus-database"`
		Port         int    `yaml:"port" env:"BANTER_BUS_DB_PORT" env-default:"27017"`
		DatabaseName string `yaml:"name" env:"BANTER_BUS_DB_NAME" env-default:"banterbus"`
		Username     string `yaml:"user" env:"BANTER_BUS_DB_USER"`
		Password     string `yaml:"password" env:"BANTER_BUS_DB_PASSWORD"`
		MaxConns     int    `yaml:"maxConns" env:"BANTER_BUS_DB_MAXCONNS" env-default:"50"`
		Timeout      int    `yaml:"timeout" env:"BANTER_BUS_DB_TIMEOUT" env-default:"3"`
	} `yaml:"database"`
}

// NewConfig creates a new config object.
func NewConfig() (config Config, err error) {
	path, exists := os.LookupEnv("BANTER_BUS_CONFIG_PATH")

	configPath := "config.yml"
	if exists {
		configPath = path
	}

	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		err = fmt.Errorf("error reading config file: %w", err)
		return Config{}, err
	}

	return config, config.validateConfig()
}

func (config *Config) validateConfig() (err error) {
	validEnvironments := map[string]bool{
		"development": true,
		"production":  true,
	}

	if !validEnvironments[config.App.Environment] {
		return fmt.Errorf("invalid environment %s", config.App.Environment)
	}

	validLogLevels := map[string]bool{
		"TRACE":   true,
		"DEBUG":   true,
		"INFO":    true,
		"WARNING": true,
		"ERROR":   true,
		"FATAL":   true,
		"PANIC":   true,
	}

	if !validLogLevels[config.App.LogLevel] {
		return fmt.Errorf("invalid log level %s", config.App.LogLevel)
	}

	if net.ParseIP(config.Webserver.Host) == nil {
		return fmt.Errorf("invalid host ip address %s", config.Webserver.Host)
	}

	if config.Webserver.Port < 1024 || config.Webserver.Port > 65535 {
		return fmt.Errorf("invalid host port %v", config.Webserver.Port)
	}

	if config.Database.Port < 1024 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port %v", config.Webserver.Port)
	}

	return err
}
