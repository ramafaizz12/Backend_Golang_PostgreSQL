package main

import (
	"nbfriends/apps/config"
	"nbfriends/apps/controller"
	"nbfriends/apps/pkg/token"
	"nbfriends/apps/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	AuthController := controller.AuthController{
		Db: db,
	}
	router.GET("/PING", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "ok",
		})
	})
	router.POST("/register", AuthController.Register)
	router.GET("/profile", CheckAuth(), AuthController.Profile)
	router.POST("/login", AuthController.Login)

	router.Run(":4444")
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.ResponseApi{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
			}
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			resp := response.ResponseApi{
				StatusCode: http.StatusUnauthorized,
				Message:    "INVALID TOKEN",
				Payload:    err.Error(),
			}
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Next()
	}
}
