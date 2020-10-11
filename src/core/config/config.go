//Package config ...
package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var once sync.Once

// Config is the data type for the expected config file.
type Config struct {
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
	var configPath = "config.yaml"
	if exists {
		configPath = path
	}

	var cfg Config
	once.Do(func() {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	})

	return &cfg
}
