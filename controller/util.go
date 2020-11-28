package controller

import (
	"admin/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, util.Response{
		Code:    1,
		Message: err.Error(),
		Data:    nil,
	})
}

func Response(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, util.Response{
		Code:    0,
		Message: "",
		Data:    data,
	})
}
