package middleware

import(
	"github.com/gin-gonic/gin"
	"intelligent-lock-srv/db"
)

func MongoConnect(context *gin.Context){
	context.Set("db", db.Client.Database(db.DataBaseName))
	context.Next()
}

