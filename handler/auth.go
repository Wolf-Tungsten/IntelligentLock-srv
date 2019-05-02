package handler

import (
	JSON "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"intelligent-lock-srv/models"
	"intelligent-lock-srv/secret"
	"io/ioutil"
	"log"
	"net/http"
)

func AuthHandler(engine *gin.Engine) {

	auth := engine.Group("/app/auth")
	auth.POST("/", code2Token)
}

func code2Token(ctx *gin.Context) {

	db := ctx.MustGet("db").(*mongo.Database)

	requestBody := struct {
		Code string `json:"code" binding:"required"`
	}{}

	_ = ctx.BindJSON(&requestBody)

	fmt.Println(requestBody)

	wxAuthUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", secret.WXAPPID, secret.WXAPPSECRET, requestBody.Code)

	fmt.Println(wxAuthUrl)

	wxResp, _ := http.Get(wxAuthUrl)
	wxRespBody, _ := ioutil.ReadAll(wxResp.Body)
	fmt.Println(string(wxRespBody))
	wxRespJSON := &struct {
		OpenId      string `json:"openid"`
		AccessToken string `json:"access_token"`
		ErrMsg      string `json:"errmsg"`
	}{}

	_ = JSON.Unmarshal(wxRespBody, &wxRespJSON)

	if wxRespJSON.ErrMsg != "" {
		ctx.JSON(http.StatusOK, struct {
			ErrMsg string `json:"errmsg"`
		}{"微信认证出错"})
		return
	}

	sessionTokenUUID, _ := uuid.NewRandom()
	sessionTokenStr := sessionTokenUUID.String()

	recordCount, _ := db.Collection("user").CountDocuments(ctx, bson.M{"openid":wxRespJSON.OpenId})

	fmt.Println(recordCount)
	switch recordCount {
	case 0:
		// 新用户注册的情况
		_, err := db.Collection("user").InsertOne(ctx, models.User{Openid: wxRespJSON.OpenId, SessionToken: sessionTokenStr})
		if err != nil {
			log.Fatal(err)
		}
	case 1:
		// 已经注册，更新sessionToken
		_, _ = db.Collection("user").UpdateOne(ctx,
			bson.M{"openid":wxRespJSON.OpenId},
			bson.M{"$set":bson.M{"sessionToken":sessionTokenStr}})
	default:
		// 完蛋了，删除可能错误的记录
		_, _ = db.Collection("user").DeleteMany(ctx, bson.M{"openid":wxRespJSON.OpenId})
		_, _ = db.Collection("user").InsertOne(ctx, models.User{ Openid: wxRespJSON.OpenId,SessionToken: sessionTokenStr})
	}

	ctx.JSON(http.StatusOK, struct{ SessionToken string `json:"sessionToken"`}{sessionTokenStr})

}
