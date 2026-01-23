package main

import (
	"log"

	"github.com/codepnw/go-starter-kit/internal/config"
	"github.com/codepnw/go-starter-kit/internal/middleware"
	"github.com/codepnw/go-starter-kit/internal/server"
	"github.com/codepnw/go-starter-kit/pkg/database"
	jwttoken "github.com/codepnw/go-starter-kit/pkg/jwt"
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

	// JWT Token
	token, err := jwttoken.NewJWTToken(cfg.JWT.AppName, cfg.JWT.SecretKey, cfg.JWT.RefreshKey)
	if err != nil {
		log.Fatal(err)
	}

	// Middleware
	mid := middleware.InitMiddleware(token)

	// Server Config
	s := server.NewServer(&server.ServerConfig{
		DB:         db,
		Token:      token,
		Middleware: mid,
	})
	// Server Run
	if err := s.SetupRouter().Run(cfg.GetAppAddress()); err != nil {
		log.Fatal(err)
	}
}
