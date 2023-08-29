package controller

import (
	"database/sql"
	"nbfriends/apps/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	Db *sql.DB
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Imgurl   string `json:"img_url"`
}

var (
	queryCreate = `
		INSERT INTO auth (email, password, img_url)
		VALUES ($1, $2, $3)
	`
)

func (a *AuthController) Register(ctx *gin.Context) {

	var req = RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Errorki": err.Error(),
		})
	}
	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Errorki": err.Error(),
		})
	}
	stat, err := a.Db.Prepare(queryCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Errorki": err.Error(),
		})
	}
	_, err = stat.Exec(
		req.Email,
		req.Password,
		req.Imgurl,
	)
	// fmt.Println(res)

	resp := response.ResponseApi{
		StatusCode: http.StatusCreated,
		Message:    "Register Sukses",
		Payload:    req,
	}
	ctx.JSON(resp.StatusCode, resp)
}
