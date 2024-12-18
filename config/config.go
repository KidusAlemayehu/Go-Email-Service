package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadENV() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
