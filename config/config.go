package config

import (
	"fmt"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// Application configuration
type AppConfig struct {
	Port     int `yaml:"port"`
	Database struct {
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
	local := "127.0.0.1"
	fmt.Println(local)
	server := "172.17.0.1"
	fmt.Println(server)
	var defaultConfig AppConfig
	defaultConfig.Port = 3000 // ganti 8080
	defaultConfig.Database.Driver = "mysql"
	defaultConfig.Database.Name = "event_db"
	defaultConfig.Database.Address = server
	defaultConfig.Database.Port = 3306
	defaultConfig.Database.Username = "root"
	defaultConfig.Database.Password = "group3"

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
