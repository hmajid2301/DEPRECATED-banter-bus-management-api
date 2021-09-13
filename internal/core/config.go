package core

import (
	"fmt"
	"net"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Conf struct {
	App struct {
		Env      string `yaml:"environment" env:"BANTER_BUS_ENVIRONMENT" env-default:"production"`
		LogLevel string `yaml:"logLevel" env:"BANTER_BUS_LOG_LEVEL" env-default:"debug"`
	} `yaml:"app"`
	Srv struct {
		Host string   `yaml:"host" env:"BANTER_BUS_WEBSERVER_HOST" env-default:"0.0.0.0"`
		Port int      `yaml:"port" env:"BANTER_BUS_WEBSERVER_PORT" env-default:"8080"`
		Cors []string `yaml:"cors" env:"BANTER_BUS_WEBSERVER_CORS"`
	} `yaml:"webserver"`
	DB struct {
		Host     string `yaml:"host" env:"BANTER_BUS_DB_HOST" env-default:"banter-bus-database"`
		Port     int    `yaml:"port" env:"BANTER_BUS_DB_PORT" env-default:"27017"`
		Name     string `yaml:"name" env:"BANTER_BUS_DB_NAME" env-default:"banterbus"`
		Username string `yaml:"user" env:"BANTER_BUS_DB_USER"`
		Password string `yaml:"password" env:"BANTER_BUS_DB_PASSWORD"`
		AuthDB   string `yaml:"authDB" env:"BANTER_BUS_AUTH_DB_NAME"`
		MaxConns int    `yaml:"maxConns" env:"BANTER_BUS_DB_MAXCONNS" env-default:"50"`
		Timeout  int    `yaml:"timeout" env:"BANTER_BUS_DB_TIMEOUT" env-default:"3"`
	} `yaml:"database"`
}

func NewConfig() (conf Conf, err error) {
	path, exists := os.LookupEnv("BANTER_BUS_CONFIG_PATH")

	confPath := "config.yml"
	if exists {
		confPath = path
	}

	err = cleanenv.ReadConfig(confPath, &conf)
	fmt.Printf("CORS %s", conf.Srv.Cors)
	if err != nil {
		err = fmt.Errorf("error reading config file %w", err)
		return Conf{}, err
	}

	return conf, conf.validate()
}

func (conf *Conf) validate() (err error) {
	validEnvs := map[string]bool{
		"development": true,
		"production":  true,
	}

	if !validEnvs[conf.App.Env] {
		return fmt.Errorf("invalid environment %s", conf.App.Env)
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

	if !validLogLevels[conf.App.LogLevel] {
		return fmt.Errorf("invalid log level %s", conf.App.LogLevel)
	}

	if net.ParseIP(conf.Srv.Host) == nil {
		return fmt.Errorf("invalid host ip address %s", conf.Srv.Host)
	}

	const maxPortRange = 65535
	const minPortRange = 1024

	if conf.Srv.Port > maxPortRange {
		return fmt.Errorf("invalid host port %v", conf.Srv.Port)
	}

	if conf.DB.Port < minPortRange || conf.DB.Port > maxPortRange {
		return fmt.Errorf("invalid database port %v", conf.Srv.Port)
	}

	return err
}
