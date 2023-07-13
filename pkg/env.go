package pkg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadEnvVar loads an environment variable from a .env file and returns it as a string.
func LoadEnvVar(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
