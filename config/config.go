package config

import "log"
import "github.com/joho/godotenv"
import "os"

func init() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Fatal(e)
	}
}

func GetConfig(key string) string {
	return os.Getenv(key)
}
