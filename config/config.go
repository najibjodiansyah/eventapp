package config

import (
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type key string

// Application configuration
type AppConfig struct {
	SecretJWT  string `yaml:"secretjwt"`
	ContextKey key    `yaml:"contextkey"`
	Port       int    `yaml:"port"`
	Database   struct {
		Driver   string `yaml:"driver"`
		Name     string `yaml:"name"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

// Initiatilize config in singleton way
func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {
	defaultConfig := AppConfig{}
	address, username, password := "127.0.0.1", "gotama", "jaladri24" // local server
	// address, username, password := "172.17.0.1", "root", "group3" // remote server

	defaultConfig.SecretJWT = "FCfUqpQSWVJN7HwT8QbCeYxdXH6JQ8pgcf9WSfM77RkzJHPHcU"
	defaultConfig.ContextKey = "EchoContextKey"
	defaultConfig.Port = 3000
	defaultConfig.Database.Driver = "mysql"
	defaultConfig.Database.Name = "event_db"
	defaultConfig.Database.Address = address
	defaultConfig.Database.Port = 3306
	defaultConfig.Database.Username = username
	defaultConfig.Database.Password = password

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		return &defaultConfig
	}

	var finalConfig AppConfig
	err := viper.Unmarshal(&finalConfig)

	if err != nil {
		log.Info("failed to extract config, will use default value")
		return &defaultConfig
	}

	return &finalConfig
}
