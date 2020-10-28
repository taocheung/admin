package main

import (
	"admin/config"
	"admin/crontab"
	"admin/model"
	"admin/router"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main()  {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	config.Init()
	router.Init(engine)
	model.Init()
	err := crontab.RemoveStatic()
	if err != nil {
		logrus.Fatal(err)
	}

	engine.Run(":8080")
}
