package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// JsonBack 是一个用于构建 JSON 响应的函数。
// 它根据传入的参数 message、ret 和 data 生成相应的 JSON 响应并返回给客户端。
// 参数:
// - c *gin.Context: Gin 框架的上下文，用于处理 HTTP 请求和响应。
// - message string: 响应消息，用于向客户端提供额外的信息。
// - ret int: 响应状态码，用于指示请求的结果状态。
// - data interface{}: 响应数据，包含客户端请求的数据结果，可以是任意类型。
func JsonBack(c *gin.Context, message string, ret int, data interface{}) {
	if ret == 0 {
		if data != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": message,
				"data":    data,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": message,
			})
		}
	} else if ret == -2 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": message,
		})
	} else if ret == -1 {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": message,
		})
	}
}
