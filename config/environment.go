package config

import (
	"os"

	log "github.com/jlentink/yaglogger"
)

func Validate() {
	validateEnvVar("DEBUG_LEVEL")
	validateEnvVar("DD_SITE")
	validateEnvVar("DD_API_KEY")
	validateEnvVar("PORT")
	validateEnvVar("BATCH_SIZE")
	validateEnvVar("BATCH_WAIT_IN_SECONDS")
}

func validateEnvVar(name string) {
	defaultValues := map[string]string{
		"DEBUG_LEVEL":           "info",
		"BATCH_SIZE":            "50",
		"BATCH_WAIT_IN_SECONDS": "5",
		"PORT":                  "1053",
	}

	value := os.Getenv(name)
	if value == "" {
		if defaultValue, found := defaultValues[name]; found {
			if defaultValue != "" {
				_ = os.Setenv(name, defaultValue)
				log.Info("Config value '%s' not set, using default value '%s'", name, defaultValue)
			}
		} else {
			log.Fatalf("Config value '%s' is missing or empty!", name)
		}
	}
}
