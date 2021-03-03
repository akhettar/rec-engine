package main

import (
	"os"

	"github.com/akhettar/rec-engine/app"
	log "github.com/sirupsen/logrus"
)

var redisURL string

func init() {
	if url, ok := os.LookupEnv("REDIS_URL"); !ok {
		redisURL = "redis://localhost:6379"
	} else {
		redisURL = url
	}
}

func main() {
	log.Info("starting the server on 3000")
	app.InitialiseApp(redisURL).Run(":3000")
}
