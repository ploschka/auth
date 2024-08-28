package main

import (
	"github.com/ploschka/auth/internal/server"
)

func main() {
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
