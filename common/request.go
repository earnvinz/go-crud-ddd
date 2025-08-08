package common

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func ReadBodyAndReset(c *gin.Context) ([]byte, error) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	// reset body ให้กลับมาอ่านซ้ำได้
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}
