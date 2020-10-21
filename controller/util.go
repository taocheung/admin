package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"data": "",
		"description": err.Error(),
		"error": err,
		"code": 1,
	})
}

func Response(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"description": "",
		"error": nil,
		"code": 0,
	})
}
