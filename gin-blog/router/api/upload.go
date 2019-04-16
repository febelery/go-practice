package api

import (
	"github.com/gin-gonic/gin"
	"learn/gin-blog/pkg/app"
	"learn/gin-blog/pkg/errors"
	"learn/gin-blog/pkg/logging"
	"learn/gin-blog/pkg/upload"
	"net/http"
)

// @Summary Import Image
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func UploadImage(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	file, image, err := ctx.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, errors.ERROR, nil)
		return
	}
	if image == nil {
		appG.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, errors.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, errors.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err = ctx.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, errors.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	})

}
