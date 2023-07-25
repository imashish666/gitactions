package server

import (
	"io"
	"net/http"
	"os"
	"www-api/config"
	_ "www-api/docs"
	"www-api/internal/logger"
	"www-api/internal/middleware"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter(config config.Config, log logger.ZapLogger) *gin.Engine {
	gin.ForceConsoleColor()
	log.Info("configs", map[string]interface{}{"config": config})
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	//create new instance of gin engine
	router := gin.New()
	router.Use(ginLogger.SetLogger())
	//endpoint /health-check for healthcheck purpose
	router.GET("/health-check", func(c *gin.Context) { c.String(http.StatusOK, "OK") })
	//swagger api docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//set deployment value from config in gin.Context
	router.Use(func(ctx *gin.Context) {
		ctx.Set("deployment", config.Deployment)
	})

	//recovery router to handle any panics
	router.Use(gin.Recovery())

	//authenticate middleware to verify all request
	router.Use(middleware.Authenticate)

	//add router groups and endpoints
	AddRoutes(router, config, log)
	//log request details

	return router
}
