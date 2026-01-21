package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/codepnw/go-starter-kit/internal/config"
	jwttoken "github.com/codepnw/go-starter-kit/pkg/jwt"
	"github.com/codepnw/go-starter-kit/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	EnvConfig *config.EnvConfig
	DB        *sql.DB
}

func Start(routerCfg *RouterConfig) error {
	cfg := routerCfg.EnvConfig
	// JWT Token
	token, err := jwttoken.NewJWTToken(cfg.JWT.AppName, cfg.JWT.SecretKey, cfg.JWT.RefreshKey)
	if err != nil {
		return err
	}
	_ = token
	
	// Gin Router
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		response.ResponseSuccess(c, http.StatusOK, "Go-Starter-Kit Running...")
	})	

	if err := r.Run(cfg.GetAppAddress()); err != nil {
		return fmt.Errorf("router start failed: %w", err)
	}
	return nil
}
