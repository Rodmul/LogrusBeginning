package main

import (
	"Notebook/logging"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.GetLogger()
	info := make(map[string]string)

	logger.Info("create router")
	router := gin.Default()

	router.POST("/create", func(ctx *gin.Context) {
		body := ctx.PostForm("body")
		title := ctx.PostForm("title")
		logger.Info(fmt.Sprintf("title: %s, body: %s", title, body))
		info[title] = body
	})

	router.GET("/", func(ctx *gin.Context) {
		for key, value := range info {
			fmt.Fprint(ctx.Writer, fmt.Sprintf("Заголовок: %s\nТекст: %s\n\n", key, value))
		}
	})

	router.Run()
}
