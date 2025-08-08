package customer

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-go/common"

	"github.com/gin-gonic/gin"
)

func IsEmailExisted(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := common.ReadBodyAndReset(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		var body struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or missing email"})
			c.Abort()
			return
		}

		idParam := c.Param("id")
		var excludeID *uint = nil
		if idParam != "" {
			id64, err := strconv.ParseUint(idParam, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
				c.Abort()
				return
			}
			id := uint(id64)
			excludeID = &id
		}

		customer, err := service.FindByEmail(body.Email, excludeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if customer != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			c.Abort()
			return
		}

		c.Next()
	}
}
