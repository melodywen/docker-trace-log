package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melodywen/docker-trace-log/app/proccess"
	router2 "github.com/melodywen/docker-trace-log/router"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	LoadConfig()
	router := gin.Default()

	if err := router2.RouterLoad(router); err != nil {
		log.Fatal("router loading err", err)
	}

	ctx := context.Background()
	if err := proccess.ProcessLoad(ctx); err != nil {
		log.Fatal("process loading err", err)
	}

	srv := &http.Server{
		Addr:    ":" + viper.GetString("web_server.port"),
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

func LoadConfig() {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("json")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
