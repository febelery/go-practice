package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"learn/gin_blog/pkg/errors"
	"net/http"
)

func BindAndValid(ctx *gin.Context, form interface{}) (int, int) {
	err := ctx.Bind(form)
	if err != nil {
		return http.StatusBadRequest, errors.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, errors.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, errors.INVALID_PARAMS
	}

	return http.StatusOK, errors.SUCCESS
}
