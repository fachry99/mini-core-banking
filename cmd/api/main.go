package main

import (
	"log"

	"github.com/fachry/mini-core-banking/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := config.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}
	defer db.Close()

	log.Println(cfg.AppName, "running on port", cfg.AppPort)
}
