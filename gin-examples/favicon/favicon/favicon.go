package favicon

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func New(path string) gin.HandlerFunc {
	path = filepath.FromSlash(path)
	if len(path) > 0 && !os.IsPathSeparator(path[0]) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(wd, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		panic("Invalid favicon path: " + path)
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(file)

	return func(ctx *gin.Context) {
		if ctx.Request.RequestURI != "/favicon.ico" {
			return
		}

		if ctx.Request.Method != "GET" && ctx.Request.Method != "HEAD" {
			status := http.StatusOK
			if ctx.Request.Method != "OPTIONS" {
				status = http.StatusMethodNotAllowed
			}
			ctx.Header("Allow", "GET,HEAD,OPTIONS")
			ctx.AbortWithStatus(status)
			return
		}

		ctx.Header("Content-Type", "image/x-icon")
		http.ServeContent(ctx.Writer, ctx.Request, "favicon.icon", info.ModTime(), reader)
		return
	}
}
