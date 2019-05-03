package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"intelligent-lock-srv/middleware"
	"intelligent-lock-srv/models"
	"net/http"
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

	userInfo := ctx.MustGet("user").(models.User)
	deviceUuid := ctx.Query("uuid")

	var deviceInfo models.Device

	db := ctx.MustGet("db").(*mongo.Database)

	err := db.Collection("device").FindOne(ctx, bson.M{ "uuid":deviceUuid }).Decode(&deviceInfo)

	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:400, Reason:"设备不存在"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:400, Reason:"设备查找过程出错"})
		return
	}

	var key models.Key

	if deviceInfo.Admin == userInfo.Id {
		// 如果是管理员的话直接提供密钥
	} else {
		// 不是管理员
		count, _ := db.Collection("access").CountDocuments(ctx, bson.M{"deviceUuid":deviceUuid, "user":userInfo.Id, "allowed":true})
		if count == 0 {
			// 没有授权
			ctx.JSON(http.StatusOK, models.Response{Success:false, Code:400, Reason:"无权访问该设备"})
			return
		}
	}

	err = db.Collection("keys").FindOneAndDelete(ctx, bson.M{"uuid":deviceUuid}).Decode(&key)
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:400, Reason:"当前设备无可用密钥，请检查设备状态"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:400, Reason:"密钥检索过程出错"})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Success:true, Code:200, Result:key})


}

func AdminDevices(ctx *gin.Context){

}

func Give(ctx *gin.Context){

}

func Fire(ctx *gin.Context){

}