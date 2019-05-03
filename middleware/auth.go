package middleware

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"intelligent-lock-srv/models"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {

	token := c.Request.Header.Get("session-token")
	db := c.MustGet("db").(*mongo.Database)
	// 如果token存在则进行权限鉴定
	if token != "" {
		var user models.User
		err := db.Collection("user").FindOne(c, bson.M{"sessionToken":token}).Decode(&user)
		if err == nil {
			c.Set("user", user)
			c.Next()
		} else {
			c.JSON(http.StatusOK, models.Response{Success:false, Code:http.StatusUnauthorized, Reason:"身份认证失效"})
			c.Abort()
		}
		// 处理请求

	} else {
		c.JSON(http.StatusOK, models.Response{Success:false, Code:http.StatusUnauthorized, Reason:"需要身份认证"})
		c.Abort()
	}
}