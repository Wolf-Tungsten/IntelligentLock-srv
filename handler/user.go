package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"intelligent-lock-srv/middleware"
	"intelligent-lock-srv/models"
	"net/http"
)

func UserHandler(engine *gin.Engine) {

	user := engine.Group("/app/user")
	user.Use(middleware.AuthMiddleware)
	user.POST("/", UpdateUserInfo)
	user.GET("/", FetchUserInfo)
}

func FetchUserInfo(ctx *gin.Context) {

	userInfo := ctx.MustGet("user").(models.User)
	fmt.Print(userInfo)
	ctx.JSON(http.StatusOK, userInfo)

}

func UpdateUserInfo(ctx *gin.Context) {

	userInfo := ctx.MustGet("user").(models.User)
	db := ctx.MustGet("db").(mongo.Database)
	requestBody := struct {
		Name string `json:"name" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
	}{}

	_ = ctx.BindJSON(&requestBody)
	_, _ = db.Collection("user").UpdateOne(ctx, bson.M{"_id":userInfo.Id}, bson.M{"name":requestBody.Name, "phoneNumber":requestBody.PhoneNumber})

	ctx.JSON(http.StatusOK, models.Response{Success:true, Code:http.StatusOK, Result:"ok"})

}