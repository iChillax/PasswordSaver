package settings

import (
	"github.com/labstack/gommon/log"

	"github.com/joho/godotenv"
)

var Evariables map[string]string

func Load_Evariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please check the .env files")
	}
	Evariables, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file, please check the .env files")
	}
	log.Info("Completed loading varibales from .env file")
}
