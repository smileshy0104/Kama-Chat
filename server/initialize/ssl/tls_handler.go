package ssl

import (
	"Kama-Chat/initialize/zlog"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"strconv"
)

// TlsHandler 返回一个gin的中间件处理器，用于强制HTTPS连接。
// 参数:
//
//	host: 服务器的主机名。
//	port: 服务器的端口号。
//
// 返回值:
//
//	gin.HandlerFunc: 一个用于gin框架的中间件函数。
func TlsHandler(host string, port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建一个secure中间件实例，配置SSL重定向选项。
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,                            // 启用SSL重定向。
			SSLHost:     host + ":" + strconv.Itoa(port), // 指定SSL主机和端口。
		})

		// 处理当前请求，应用安全中间件。
		err := secureMiddleware.Process(c.Writer, c.Request)

		// 如果处理过程中出现错误，记录错误并终止执行。
		if err != nil {
			zlog.Fatal(err.Error())
			return
		}

		// 继续执行后续的中间件或处理函数。
		c.Next()
	}
}
