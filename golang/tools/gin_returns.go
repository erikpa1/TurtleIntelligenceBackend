package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AutoReturn(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func Return405WithError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}
