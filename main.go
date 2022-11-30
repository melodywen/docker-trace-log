package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/melodywen/docker-trace-log/app/proccess"
	router2 "github.com/melodywen/docker-trace-log/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	router := gin.Default()

	if err := router2.RouterLoad(router); err != nil {
		log.Fatal("router loading err", err)
	}

	ctx := context.Background()
	if err := proccess.ProcessLoad(ctx); err != nil {
		log.Fatal("process loading err", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
