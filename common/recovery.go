package common

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func JSONRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		fmt.Printf("Panic recovered: %v\n", recovered)
		debug.PrintStack()

		c.JSON(http.StatusInternalServerError, ResponseError{
			Error: fmt.Sprintf("Internal Server Error: %v", recovered),
		})
		c.Abort()
	})
}
