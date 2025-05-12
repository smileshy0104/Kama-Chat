package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 创建并返回一个Gin的中间件处理器，用于处理跨域请求(CORS)。
// 该中间件配置允许来自任何源的请求，支持GET、POST、PUT、DELETE和OPTIONS方法，
// 并允许特定的HTTP头部，包括Origin、Content-Length、Content-Type和Authorization。
func Cors() gin.HandlerFunc {
	// 使用默认的跨域资源共享（CORS）配置创建一个cors.Config实例。
	corsConfig := cors.DefaultConfig()
	// 设置允许跨域请求的源列表，"*"表示允许所有源。
	corsConfig.AllowOrigins = []string{"*"}
	// 指定允许的HTTP方法列表。
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// 指定允许的HTTP头部列表，这些头部在跨域请求中可以被包含。
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	// 使用配置好的cors.Config实例创建并返回一个gin.HandlerFunc中间件。
	return cors.New(corsConfig)
}

// Cors 生成跨域请求的中间件
// 该中间件主要用于设置响应头，以允许跨域请求
// 参数: 无
// 返回值: gin.HandlerFunc 中间件函数
func CorsNew() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求方法
		method := c.Request.Method
		// 获取请求头中的 Origin 字段
		origin := c.Request.Header.Get("Origin")
		// 设置允许跨域请求的来源
		c.Header("Access-Control-Allow-Origin", origin)
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		// 设置允许暴露的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		// 设置是否允许发送 cookies
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
