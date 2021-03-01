package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/akhettar/rec-engine/app"
)

func main() {
	
	log.Info("starting the server on 3000")
	app.InitialiseApp("redis://34.66.203.46:6379").Run(":3000")
}
