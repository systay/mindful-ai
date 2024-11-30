package utils

import (
	"errors"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// Look for .env file in parent directories up to 3 levels up
	envFiles := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			return nil
		}
	}

	return errors.New(".env file not found")
}
