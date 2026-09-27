package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/config"
)

func initConfig() error {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	config.InitConfigService(config.NewConfigurationService())

	return nil
}
