package healthcheck

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {

	handler := NewHandler()
	handler.RegisterRoutes(rg)
}
