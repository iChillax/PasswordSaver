package settings

import (
	"os"

	"github.com/labstack/gommon/log"

	"github.com/joho/godotenv"
)

var Evariables map[string]string

func Load_Evariables() {
	// First, try to load from .env file (this won't override existing env vars)
	err := godotenv.Load()
	if err != nil {
		log.Warn("No .env file found, using environment variables")
	}

	// Read all variables from .env file
	envFileVars, err := godotenv.Read()
	if err != nil {
		log.Warn("Error reading .env file, will use environment variables only")
		envFileVars = make(map[string]string)
	}

	// Initialize Evariables map
	Evariables = make(map[string]string)

	// First, add all variables from .env file
	for key, value := range envFileVars {
		Evariables[key] = value
	}

	// Then, override with any environment variables that are actually set
	// This gives priority to exported environment variables
	for key := range envFileVars {
		if envValue, exists := os.LookupEnv(key); exists {
			Evariables[key] = envValue
		}
	}

	log.Info("Completed loading variables from .env file and environment")
}
