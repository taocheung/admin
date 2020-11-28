package router

import (
	"admin/controller"
	_ "admin/docs"
	"admin/middleware"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init(router *gin.Engine) {
	router.Use(middleware.Cors())
	user := router.Group("/user")
	user.POST("login", controller.Login)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
