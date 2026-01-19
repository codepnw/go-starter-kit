package main

import (
	"log"

	"github.com/codepnw/go-starter-kit/internal/config"
	"github.com/codepnw/go-starter-kit/internal/router"
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
	
	// Start Server
	err = router.Start(&router.RouterConfig{
		EnvConfig: cfg,
		DB:        db,
	})
	if err != nil {
		log.Fatal(err)
	}
}
