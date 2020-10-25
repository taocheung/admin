package main

import (
	"admin/config"
	"admin/model"
	"admin/router"
	"github.com/gin-gonic/gin"
)

func main()  {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	config.Init()
	router.Init(engine)
	model.Init()

	engine.Run(":31000")
}
