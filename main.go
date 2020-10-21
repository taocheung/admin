package main

import (
	"admin/model"
	"admin/router"
	"github.com/gin-gonic/gin"
)

func main()  {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	router.Init(engine)
	model.Init()

	engine.Run(":8080")
}
