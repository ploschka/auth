package main

import (
	log "github.com/ploschka/auth/internal/logger"
	"github.com/ploschka/auth/internal/server"
)

func main() {
	log.Info("Auth server started")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
