package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/melodywen/docker-trace-log/contracts"
	"strings"
	"time"
)

func LogMiddleWare(app contracts.AppAttributeInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		id := strings.Replace(uuid.New().String(), "-", "", -1)
		c.Set("trace_id", id)
		// 请求前
		c.Next()
		// 请求后
		latency := time.Since(t)
		app.GetLog().Debug(c, "请求耗时为：%v", latency)
	}
}
