//Package config ...
package config

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ilyakaznacheev/cleanenv"
)

var once sync.Once

// Config is the data type for the expected config file.
type Config struct {
	App struct {
		Environment string `yaml:"environment" env:"BANTER_BUS_ENVIRONMENT" env-default:"production"`
		LogLevel    string `yaml:"logLevel" env:"BANTER_BUS_LOG_LEVEL" env-default:"debug"`
	} `yaml:"app"`
	Database struct {
		Host         string `yaml:"host" env:"BANTER_BUS_DB_HOST" env-default:"banter-bus-database"`
		Port         string `yaml:"port" env:"BANTER_BUS_DB_PORT" env-default:"27017"`
		DatabaseName string `yaml:"name" env:"BANTER_BUS_DB_NAME" env-default:"banterbus"`
		Username     string `yaml:"user" env:"BANTER_BUS_DB_USER"`
		Password     string `yaml:"password" env:"BANTER_BUS_DB_PASSWORD"`
	} `yaml:"database"`
}

//GetConfig ...
func GetConfig() *Config {
	path, exists := os.LookupEnv("BANTER_BUS_CONFIG_PATH")
	var configPath = "config.yml"
	if exists {
		configPath = path
	}

	var cfg Config
	once.Do(func() {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Error("Failed to load config.", err)
			os.Exit(1)
		}
	})

	return &cfg
}
