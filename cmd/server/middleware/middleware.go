package middleware

import (
	"os"
	"time"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func ValidarToken() gin.HandlerFunc {
	//Verificacion de Token
	return func(ctx *gin.Context) {
		userToken := ctx.GetHeader("TOKEN")
		token := os.Getenv("TOKEN")
		if userToken != token {
			ctx.AbortWithStatusJSON(401, web.InvalidTokenResponse)
			return
		}
		ctx.Next()
	}
}

func Recovery() gin.HandlerFunc {
	//Verificacion de Token
	return func(ctx *gin.Context) {

		defer func() {
			_ = recover()

			var time time.Time = time.Now()

			bytes, _ := json.Marshal(gin.H{
				"verb": ctx.Request.Method,
				"date": time.String(),
				"url":  ctx.Request.RequestURI,
				"size": ctx.Request.ContentLength,
			})

			gin.DefaultErrorWriter.Write(bytes)
			gin.DefaultErrorWriter.Write([]byte("\n"))
			ctx.AbortWithStatusJSON(500, web.InternalErrorResponse)
		}()

		ctx.Next()
	}
}
