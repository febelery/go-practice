package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"learn/gin-blog/pkg/errors"
	"learn/gin-blog/pkg/util"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		var data interface{}

		code = errors.SUCCESS
		token := ctx.Query("token")

		if token == "" {
			code = errors.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = errors.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = errors.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != errors.SUCCESS {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  errors.GetMsg(code),
				"data": data,
			})

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
