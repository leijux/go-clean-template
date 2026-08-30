package main

import (
	"log"

	"github.com/leijux/go-clean-template/config"
	"github.com/leijux/go-clean-template/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
