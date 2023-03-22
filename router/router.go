package router

import (
	"github.com/gin-gonic/gin"
	"github.com/melodywen/docker-trace-log/app/http/middleware"
	"github.com/melodywen/docker-trace-log/contracts"
	"net/http"
)

// RouteLoad 加载路由
//  @param router
//  @return error
func RouteLoad(router *gin.Engine, app contracts.AppAttributeInterface) error {

	router.Use(middleware.LogMiddleWare(app))

	router.Static("/public", "./public")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/public/index.html")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return nil
}
