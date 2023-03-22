package provider

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/melodywen/docker-trace-log/router"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type HttpProvider struct {
	server *http.Server
}

func NewHttpProvider() *HttpProvider {
	return &HttpProvider{}
}

func (h *HttpProvider) StartServerBeforeEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()

	return nil
}

func (h *HttpProvider) StartServerAfterEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()

	route := gin.Default()

	if err := router.RouteLoad(route, app); err != nil {
		app.GetLog().Panic(ctx, "router loading err", err)
	}

	h.server = &http.Server{
		Addr:    ":" + app.GetConfig().GetString("web_server.port"),
		Handler: route,
	}
	go func() {
		// 服务连接
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.GetLog().Panic(ctx, "listen: %s\n", err)
		}
	}()

	h.Stop(ctx, app)
	return nil
}

func (h *HttpProvider) Stop(ctx context.Context, app contracts.AppAttributeInterface) {

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.GetLog().Info(ctx, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		app.GetLog().Panic(ctx, "Server Shutdown:", err)
	}

	app.GetLog().Info(ctx, "Server exiting")
}
