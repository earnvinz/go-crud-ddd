package common

import "github.com/gin-gonic/gin"

type ResponseError struct {
	Error string `json:"error" example:"error message"`
}

func JSONError(c *gin.Context, status int, errMsg string) {
	c.JSON(status, ResponseError{Error: errMsg})
}
