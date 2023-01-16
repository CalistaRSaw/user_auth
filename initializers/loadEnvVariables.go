package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() { //first letter capital so other packages can use
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
