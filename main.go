package main

import (
	"admin/config"
	"admin/model"
	"admin/router"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Admin System API
// @version 1.0
// @description This is a admin server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v1
// @host 127.0.0.1:8080
func main() {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	config.Init()
	router.Init(engine)
	model.Init()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatal(err)
		}
	}()

	quite := make(chan os.Signal)
	signal.Notify(quite, os.Interrupt)
	select {
	case <-quite:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logrus.Fatal(err)
		}
	}
}
