package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	basepath := "gin-examples/upload_file/multiple/"
	router.Static("/", basepath+"public")
	router.POST("/upload", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")

		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		files := form.File["files"]

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			if err := ctx.SaveUploadedFile(file, basepath+filename); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
		}
		ctx.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email))
	})

	router.Run(":8080")
}
