package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	router := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20

	basepath := "gin-examples/upload_file/single/"

	router.Static("/", basepath+"public/")
	router.POST("/upload", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")

		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := ctx.SaveUploadedFile(file, basepath+filename); err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		ctx.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", filename, name, email))
	})

	router.Run(":8080")
}
