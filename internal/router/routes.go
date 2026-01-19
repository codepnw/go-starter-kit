package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/codepnw/go-starter-kit/internal/config"
	"github.com/codepnw/go-starter-kit/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	EnvConfig *config.EnvConfig
	DB        *sql.DB
}

func Start(routerCfg *RouterConfig) error {
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		response.ResponseSuccess(c, http.StatusOK, "Go-Starter-Kit Running...")
	})

	if err := r.Run(routerCfg.EnvConfig.GetAppAddress()); err != nil {
		return fmt.Errorf("router start failed: %w", err)
	}
	return nil
}
