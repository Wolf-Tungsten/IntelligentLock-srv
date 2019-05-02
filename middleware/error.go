package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorWrapper(context *gin.Context) {
	context.Next()

	if len(context.Errors) > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error":context.Errors})
	}

}