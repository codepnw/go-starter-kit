package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/codepnw/go-starter-kit/internal/middleware"
	jwttoken "github.com/codepnw/go-starter-kit/pkg/jwt"
	"github.com/codepnw/go-starter-kit/pkg/utils/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	DB         *sql.DB
	Token      *jwttoken.JWTToken
	Middleware *middleware.Middleware
}

type server struct {
	cfg *ServerConfig
}

func NewServer(cfg *ServerConfig) *server {
	return &server{cfg: cfg}
}

func (s *server) SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(s.cfg.Middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		response.ResponseSuccess(c, http.StatusOK, "Go-Starter-Kit Running...")
	})

	s.registerUserRoutes()

	return r
}

func (s *server) registerUserRoutes() {}
