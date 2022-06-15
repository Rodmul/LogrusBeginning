package main

import (
	"Notebook/logging"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		time.Sleep(2 * time.Second)
		for key, value := range info {
			fmt.Fprint(ctx.Writer, fmt.Sprintf("Title: %s\nBody: %s\n\n", key, value))
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to graceful shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ... ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Serer Shutdown: %v", err)
	}
	// catching ctx.Done().timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Info("timeout 5 seconds.")
	}
	logger.Info("Server exiting")
}
