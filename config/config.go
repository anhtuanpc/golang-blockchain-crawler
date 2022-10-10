package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Fatal(e)
	}
}

func GetConfig(key string) string {
	return os.Getenv(key)
}
