package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"learn/gin-blog/pkg/app"
	"learn/gin-blog/pkg/errors"
	"learn/gin-blog/pkg/util"
	"learn/gin-blog/service/auth_service"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	valid := validation.Validation{}

	username := ctx.Query("username")
	password := ctx.Query("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, errors.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, map[string]string{
		"token": token,
	})

}
