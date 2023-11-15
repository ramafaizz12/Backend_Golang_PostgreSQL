package controller

import (
	"database/sql"
	"nbfriends/apps/pkg/token"
	"nbfriends/apps/response"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Db *sql.DB
}

type RegisterRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
}

type LoginRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
}

type Auth struct {
	Id           int
	Email        string
	password     string
	Organization string
}

var (
	queryCreate = `
		INSERT INTO auth (email, password, organization)
		VALUES ($1, $2, $3)
	`

	queryFindByEmail = `
		SELECT id, email, password, organization
		FROM auth
		WHERE email=$1
	`
)

func (a *AuthController) Register(ctx *gin.Context) {

	var req = RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Errorki": err.Error(),
		})
		return
	}
	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Errorki": err.Error(),
		})
		return
	}
	EmailValid := isEmailValid(req.Email)
	if !EmailValid {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Errorki": "Format Email Salah",
		})
		return
	}
	if req.Organization != "farmer" && req.Organization != "distributor" && req.Organization != "retailer" && req.Organization != "wholesaler" && req.Organization != "consumer" && req.Organization != "manufacturer" {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"Error": "Tolong Masukkan Organisasi Anda",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Errorki": err.Error(),
		})
		return
	}
	req.Password = string(hash)

	stat, err := a.Db.Prepare(queryCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Errorki": err.Error(),
		})
		return
	}
	_, err = stat.Exec(
		req.Email,
		req.Password,
		req.Organization,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Errorki": err.Error(),
		})
		return
	}

	resp := response.ResponseApi{
		StatusCode: http.StatusCreated,
		Message:    "Register Sukses",
	}
	ctx.JSON(resp.StatusCode, resp)
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req = LoginRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	stmt, err := a.Db.Prepare(queryFindByEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
	row := stmt.QueryRow(req.Email)

	var auth = Auth{}

	err = row.Scan(
		&auth.Id,
		&auth.Email,
		&auth.password,
		&auth.Organization,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(auth.password), []byte(req.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Error": err.Error(),
		})
		return
	}
	if auth.Organization != req.Organization {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"Error": "Nama Organisasi Salah",
		})
		return
	}

	tok := token.PayloadToken{
		AuthId: auth.Id,
	}
	tokString, err := token.GenerateToken(&tok)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	resp := response.ResponseApi{
		StatusCode: http.StatusOK,
		Message:    "Login Sukses",
		Payload: gin.H{
			"token": tokString,
		},
	}
	println(auth.Organization)
	ctx.JSON(resp.StatusCode, resp)
}

func (a *AuthController) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"id": ctx.GetInt("authId"),
	})
}
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	kondisi := emailRegex.MatchString(e)

	return kondisi
}
