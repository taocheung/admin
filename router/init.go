package router

import (
	"admin/controller"
	"admin/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	user := router.Group("/user")
	user.POST("login", controller.Login)
	user.Use(middleware.JWTUserAuth())
	{
		user.POST("add", controller.AddUser)
		user.POST("update", controller.UpdateUser)
		user.POST("delete", controller.DeleteUser)
		user.POST("list", controller.ListUser)
	}

	resource := router.Group("/resource")
	resource.Use(middleware.JWTAuth())
	{
		resource.POST("import", controller.ResourceImport)
		resource.GET("export", controller.ResourceExport)
		resource.POST("list", controller.ResourceList)
	}
	router.GET("download", controller.Template)

}
