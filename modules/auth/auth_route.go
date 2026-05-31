package auth

import (
	moduleLib "margin-delver/lib"
	authprovider "margin-delver/modules/auth/auth_provider"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouterGroupFunc func(router *gin.RouterGroup)

func RouterGroupWithDeps(
	log *moduleLib.BaseLog,
	cfg *moduleLib.AppConfig,
	db *gorm.DB,
) RouterGroupFunc {
	return func(router *gin.RouterGroup) {
		handler := authprovider.InitializeAuthHandler(log, cfg, db)

		router.POST("/auth/login", handler.Login)
	}
}
