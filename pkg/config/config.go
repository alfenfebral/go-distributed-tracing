package config

import (
	"errors"
	"go-distributed-tracing/utils"

	"github.com/joho/godotenv"
)

// LoadConfig - load environment config
func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		utils.CaptureError(errors.New("error loading .env file"))
	}

	return err
}
