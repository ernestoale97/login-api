package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	AppDatabase string `mapstructure:"APP_DATABASE"`
	AppPort     string `mapstructure:"APP_PORT"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./env")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&configuration)

	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		return
	}

	return configuration, nil
}
