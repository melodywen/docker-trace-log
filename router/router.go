package router

import "github.com/gin-gonic/gin"

// RouterLoad 加载路由
//  @param router
//  @return error
func RouterLoad(router *gin.Engine) error {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return nil
}
