package v1

import (
	moduleLib "margin-delver/lib"
	auth "margin-delver/modules/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetRouter(
	router *gin.Engine,
	cfg *moduleLib.AppConfig,
	log *moduleLib.BaseLog,
	db *gorm.DB,
) {
	internalGroup := router.Group("/internal/v1")
	auth.RouterGroupWithDeps(log, cfg, db)(internalGroup)
}
