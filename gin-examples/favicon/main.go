package main

import (
	"github.com/gin-gonic/gin"
	"learn/gin-examples/favicon/favicon"
	"net/http"
)

func main() {
	app := gin.Default()

	app.Use(favicon.New("gin-examples/favicon/favicon.ico"))

	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello favicon.")
	})

	app.Run(":8080")
}
