package main

import (
	"nbfriends/apps/config"
	"nbfriends/apps/controller"
	"net/http"

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

	router.Run(":4444")
}
