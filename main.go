package main

import (
	"dd-log-proxy/config"
	"dd-log-proxy/server"
	"os"
	"strings"

	log "github.com/jlentink/yaglogger"
)

var (
	version = "dev"
)

func main() {
	log.Info("Starting dd-log-proxy version %s", version)

	switch strings.ToLower(os.Getenv("DEBUG_LEVEL")) {
	case "info":
		log.SetLevel(log.LevelInfo)
	case "debug":
		log.SetLevel(log.LevelDebug)
	case "warning":
		log.SetLevel(log.LevelWarn)
	case "fatal":
		log.SetLevel(log.LevelFatal)
	default:
		log.SetLevel(log.LevelInfo)
	}
	log.Debug("Validating environment..")
	config.Validate()

	server.Start()

}
