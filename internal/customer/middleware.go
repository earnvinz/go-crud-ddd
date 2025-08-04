package customer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsEmailExisted(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

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

		var body struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or missing email"})
			c.Abort()
			return
		}

		customer, err := repo.FindByEmail(body.Email, excludeID)
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
