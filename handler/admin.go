package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	db := ctx.MustGet("db").(*mongo.Database)
	userInfo := ctx.MustGet("user").(models.User)

	type accessItem struct {
		Id primitive.ObjectID `json:"userId"`
		Name string `json:"userName"`
		PhoneNumber string `json:"phoneNumber"`
		Allowed bool `json:"allowed"`
	}

	type deviceItem struct {
		Uuid string `json:"uuid" bson:"uuid"`
		AccessList []accessItem `json:"accessList,omitempty"`
	}

	var deviceList []deviceItem

	deviceCursor, _ := db.Collection("device").Find(ctx, bson.M{"admin":userInfo.Id})

	for deviceCursor.Next(ctx) {
		var _device models.Device
		_ = deviceCursor.Decode(&_device)

		device := deviceItem{Uuid:_device.Uuid}

		accessCursor, _ := db.Collection("access").Find(ctx, bson.M{"deviceUuid":device.Uuid})
		for accessCursor.Next(ctx) {
			var _access models.DeviceAccess
			_ = accessCursor.Decode(&_access)
			var _accessUser models.User
			_ = db.Collection("user").FindOne(ctx, bson.M{"_id": _access.User}).Decode(&_accessUser)
			fmt.Println(_accessUser)
			device.AccessList = append(device.AccessList, accessItem{Id:_accessUser.Id, Name:_accessUser.Name, PhoneNumber:_accessUser.PhoneNumber, Allowed:_access.Allowed})
		}
		_ = accessCursor.Close(ctx)
		deviceList = append(deviceList, device)
	}
	_ = deviceCursor.Close(ctx)

	ctx.JSON(http.StatusOK, models.Response{Success:true, Code:200, Result:deviceList})
}

func Give(ctx *gin.Context){

	db := ctx.MustGet("db").(*mongo.Database)
	userInfo := ctx.MustGet("user").(models.User)

	requestBody := struct {
		DeviceUuid string `json:"deviceUuid"`
		UserId string `json:"userId"`
	}{}

	_ = ctx.BindJSON(&requestBody)

	var device models.Device

	err := db.Collection("device").FindOne(ctx, bson.M{"uuid":requestBody.DeviceUuid}).Decode(&device)

	if err == mongo.ErrNoDocuments {
		ctx.JSON(200, models.Response{Success:false, Code:400, Reason:"设备不存在"})
		return
	} else if err != nil {
		ctx.JSON(200, models.Response{Success:false, Code:400, Reason:"设备检索出错"})
		return
	}

	if device.Admin != userInfo.Id {
		ctx.JSON(200, models.Response{Success:false, Code:400, Reason:"无权操作"})
		return
	}

	userId, _ := primitive.ObjectIDFromHex(requestBody.UserId)
	_, _ = db.Collection("access").UpdateMany(ctx, bson.M{"deviceUuid":requestBody.DeviceUuid, "user":userId}, bson.M{"$set":bson.M{"allowed":true}})
	ctx.JSON(200, models.Response{Success:true, Code:200, Reason:"授权成功"})
}

func Fire(ctx *gin.Context){

	db := ctx.MustGet("db").(*mongo.Database)
	userInfo := ctx.MustGet("user").(models.User)


	DeviceUuid := ctx.Query("deviceUuid")
	UserId := ctx.Query("userId")


	var device models.Device

	err := db.Collection("device").FindOne(ctx, bson.M{"uuid":DeviceUuid}).Decode(&device)

	if err == mongo.ErrNoDocuments {
		ctx.JSON(200, models.Response{Success:false, Code:400, Reason:"设备不存在"})
		return
	} else if err != nil {
		ctx.JSON(200, models.Response{Success:false, Code:400, Reason:"设备检索出错"})
		return
	}

	if device.Admin != userInfo.Id {
		ctx.JSON(200, models.Response{Success:false, Code:400, Reason:"无权操作"})
		return
	}

	userId, _ := primitive.ObjectIDFromHex(UserId)
	_, _ = db.Collection("access").DeleteMany(ctx, bson.M{"deviceUuid":DeviceUuid, "user":userId})
	ctx.JSON(200, models.Response{Success:true, Code:200, Reason:"撤销成功"})

}