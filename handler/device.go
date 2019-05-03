package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"intelligent-lock-srv/models"
	"net/http"
)

func DeviceHandler(engine *gin.Engine) {

	device := engine.Group("/device")
	device.POST("/activate", Activate)
	device.POST("/uploadkeys", UploadKeys)
}

func Activate(ctx *gin.Context) {
	db := ctx.MustGet("db").(*mongo.Database)

	requestBody := struct {
		Uuid string `json:"uuid"`
		UserSession string `json:"userSessionToken"`
	}{}

	_ = ctx.BindJSON(&requestBody)

	// 获取session
	var userInfo models.User
	err := db.Collection("user").FindOne(ctx, bson.M{"sessionToken":requestBody.UserSession}).Decode(&userInfo)

	if err != nil {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:http.StatusBadRequest, Reason:"用户会话无效"})
		return
	}

	// 根据Uuid查询锁的状态
	var deviceInfo models.Device
	err = db.Collection("device").FindOne(ctx, bson.M{"uuid":requestBody.Uuid}).Decode(&deviceInfo)

	fmt.Print(err)
	if err == mongo.ErrNoDocuments {
		// 处于可激活状态
		_, _ = db.Collection("device").InsertOne(ctx, models.Device{Uuid:requestBody.Uuid, Admin:userInfo.Id})
		ctx.JSON(http.StatusOK, models.Response{Success:true, Code:http.StatusOK, Result:"激活成功"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:http.StatusBadRequest, Reason:"设备检索出错"})
		return
	}

	// 执行到此处说明设备存在并已经被激活，那么激活请求被视为申请访问
	if deviceInfo.Admin != userInfo.Id {
		_, err = db.Collection("access").InsertOne(ctx, models.DeviceAccess{DeviceUuid:requestBody.Uuid, User:userInfo.Id, Allowed:false})
		ctx.JSON(http.StatusOK, models.Response{Success:true, Code:http.StatusOK, Result:"申请提交成功"})
	} else {
		ctx.JSON(http.StatusOK, models.Response{Success:false, Code:http.StatusBadRequest, Reason:"重复激活"})
	}

}

func UploadKeys(ctx *gin.Context) {

	db := ctx.MustGet("db").(*mongo.Database)

	requestBody := struct {
		Uuid string `json:"uuid"`
		Keys []string `json:"keys"`
	}{}

	_ = ctx.BindJSON(&requestBody)

	var err error
	for _, key := range requestBody.Keys {
		_, err = db.Collection("keys").InsertOne(ctx, models.Key{Uuid:requestBody.Uuid, Key:key})
		if err != nil {
			ctx.JSON(http.StatusOK, models.Response{Success:false, Code:http.StatusBadRequest, Reason:"密钥插入出错"})
			return
		}
	}

	ctx.JSON(http.StatusOK, models.Response{Success:true, Code:http.StatusOK, Result:fmt.Sprintf("成功插入%d个密钥", len(requestBody.Keys))})

}
