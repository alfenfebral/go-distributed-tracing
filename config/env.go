package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig - load environment config
func LoadConfig() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	return nil
}
