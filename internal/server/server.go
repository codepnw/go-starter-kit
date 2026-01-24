package server

import (
	"database/sql"
	"net/http"
	"time"

	userhandler "github.com/codepnw/go-starter-kit/internal/features/user/handler"
	userrepository "github.com/codepnw/go-starter-kit/internal/features/user/repository"
	userservice "github.com/codepnw/go-starter-kit/internal/features/user/service"
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

type Server struct {
	cfg    *ServerConfig
	router *gin.Engine
}

func NewServer(cfg *ServerConfig) *Server {
	r := gin.New()

	s := &Server{
		cfg:    cfg,
		router: r,
	}

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

	return s
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) registerUserRoutes() {
	repo := userrepository.NewUserRepository(s.cfg.DB)
	service := userservice.NewUserService(s.cfg.Token, repo)
	handler := userhandler.NewUserHandler(service)

	users := s.router.Group("/users")

	users.POST("/register", handler.Register)
	users.POST("/login", handler.Login)
}
