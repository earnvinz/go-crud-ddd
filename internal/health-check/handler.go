package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// HealthHandler godoc
// @BasePath /api/v1
// @Summary      Health Check
// @Description  Returns OK
// @Tags         Health
// @Success      200 {string} string "OK"
// @Router       /health-check [get]
func (h *Handler) HealthCheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	healthCheck := rg.Group("/health-check")

	healthCheck.GET("/", h.HealthCheckHandler)
}
