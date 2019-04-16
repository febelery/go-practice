package v1

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"learn/gin-blog/pkg/app"
	"learn/gin-blog/pkg/errors"
	"learn/gin-blog/pkg/export"
	"learn/gin-blog/pkg/logging"
	"learn/gin-blog/pkg/setting"
	"learn/gin-blog/pkg/util"
	"learn/gin-blog/service/tag_service"
	"net/http"
)

func GetTags(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	name := ctx.Query("name")
	state := -1
	if arg := ctx.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(ctx),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

func AddTag(ctx *gin.Context) {
	var (
		appG = app.Gin{Ctx: ctx}
		form AddTagForm
	)

	httpCode, errCode := app.BindAndValid(ctx, &form)
	if errCode != errors.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, errors.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"required:Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

func EditTag(ctx *gin.Context) {
	var (
		appG = app.Gin{Ctx: ctx}
		form = EditTagForm{ID: com.StrTo(ctx.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(ctx, &form)
	if errCode != errors.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, errors.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, nil)
}

func DeleteTag(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	valid := validation.Validation{}
	id := com.StrTo(ctx.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, errors.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, errors.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, nil)
}

func ExportTag(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	name := ctx.PostForm("name")
	state := -1
	if arg := ctx.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errors.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})

}

func ImportTag(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, errors.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, errors.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, errors.SUCCESS, nil)
}
