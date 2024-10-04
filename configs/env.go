package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	godotenv.Load()
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
