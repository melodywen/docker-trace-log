package main

import (
	"context"
	"github.com/gin-gonic/gin"
	app2 "github.com/melodywen/docker-trace-log/app"
	"github.com/melodywen/docker-trace-log/package/logs"
	"github.com/melodywen/docker-trace-log/router"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	app := app2.GetApp()
	ctx := logs.GetContextOfLog()

	app.NotifyStartServerBeforeEvent(ctx)

	route := gin.Default()

	if err := router.RouterLoad(route); err != nil {
		app.Log.Panic(ctx, "router loading err", err)
	}

	srv := &http.Server{
		Addr:    ":" + app.Config.GetString("web_server.port"),
		Handler: route,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Log.Panic(ctx, "listen: %s\n", err)
		}
	}()

	app.NotifyStartServerAfterEvent(ctx)

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.Log.Info(ctx, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		app.Log.Panic(ctx, "Server Shutdown:", err)
	}

	app.Log.Info(ctx, "Server exiting")
}
