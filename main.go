package main

import (
	"github.com/gin-gonic/gin"
	"intelligent-lock-srv/db"
	"intelligent-lock-srv/handler"
	"intelligent-lock-srv/middleware"
)

const Port = "3001"
func main() {

	db.Connect()
	r := gin.Default()

	r.Use(middleware.ErrorWrapper)
	r.Use(middleware.MongoConnect)
	r.Use(middleware.Cors)

	handler.AuthHandler(r)
	handler.UserHandler(r)
	handler.DeviceHandler(r)
	handler.AdminHandler(r)
	r.Run(":" + Port)
}
