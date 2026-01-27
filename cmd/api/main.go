package main

import (
	"log"

	"github.com/codepnw/go-starter-kit/internal/config"
	"github.com/codepnw/go-starter-kit/internal/server"
	"github.com/codepnw/go-starter-kit/pkg/database"
)

const envPath = ".env"

func main() {
	// Load Config
	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	// Connect Database
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s, err := server.NewServer(cfg, db)
	if err != nil {
		log.Fatal(err)
	}
	// Server Start
	if err := s.Start(cfg.GetAppAddress()); err != nil {
		log.Fatal(err)
	}
}
