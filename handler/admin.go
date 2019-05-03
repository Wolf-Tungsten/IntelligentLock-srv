package handler

import (
	"github.com/gin-gonic/gin"
	"intelligent-lock-srv/middleware"
)

func AdminHandler(engine *gin.Engine) {

	admin := engine.Group("/app/admin")
	admin.Use(middleware.AuthMiddleware)
	admin.GET("/key", GetKey)
	admin.GET("/", AdminDevices)
	admin.POST("/", Give)
	admin.DELETE("/", Fire)

}

func GetKey(ctx *gin.Context){

}

func AdminDevices(ctx *gin.Context){

}

func Give(ctx *gin.Context){

}

func Fire(ctx *gin.Context){

}