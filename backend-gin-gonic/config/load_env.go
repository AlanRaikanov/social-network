package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvs() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
